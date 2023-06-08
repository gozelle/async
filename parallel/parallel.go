package parallel

import (
	"context"
	"fmt"
	"sync"

	"github.com/gozelle/async"
)

type Null = async.Null

type Result[T any] struct {
	Error error
	Value T
}

type Runner[T any] async.Runner[T]

func Run[T any](ctx context.Context, limit uint, runners []Runner[T]) <-chan *Result[T] {

	results := make(chan *Result[T], len(runners))

	if limit == 0 {
		defer func() {
			results <- &Result[T]{Error: fmt.Errorf("limit expect great than 0")}
			close(results)
		}()
		return results
	}

	if ctx == nil {
		ctx = context.Background()
	}

	cctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	sem := make(chan struct{}, limit)

	for _, v := range runners {
		wg.Add(1)
		sem <- struct{}{}
		go func(runner Runner[T]) {
			defer func() {
				<-sem
				wg.Done()
			}()
			select {
			case <-cctx.Done():
				return
			default:
				r, err := runner(ctx)
				if err != nil {
					results <- &Result[T]{Error: err}
					cancel()
				} else {
					results <- &Result[T]{Value: r}
				}
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(results)
		close(sem)
		cancel()
	}()

	return results
}

func Wait[T any](results <-chan *Result[T], handler func(v T) error) error {
	for item := range results {
		if item.Error != nil {
			return item.Error
		}
		err := handler(item.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
