package parallel

import (
	"context"
	"fmt"
	"github.com/gozelle/atomic"
	"runtime/debug"
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
	
	err := atomic.NewError(nil)
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, limit)
	done := make(chan struct{})
	
	go func() {
		select {
		case <-ctx.Done():
			err.Store(ctx.Err())
		case <-done:
			return
		}
	}()
	
	for _, v := range runners {
		sem <- struct{}{}
		if err.Load() != nil {
			<-sem
			continue
		}
		wg.Add(1)
		go func(runner Runner[T]) {
			defer func() {
				e := recover()
				if e != nil {
					err.Store(fmt.Errorf("%v", e))
					debug.PrintStack()
				}
				<-sem
				wg.Done()
			}()
			
			r, e := runner(ctx)
			if e != nil {
				err.Store(e)
			} else {
				results <- &Result[T]{Value: r}
			}
		}(v)
	}
	
	go func() {
		wg.Wait()
		if err.Load() != nil {
			results <- &Result[T]{Error: err.Load()}
		}
		close(done)
		close(results)
		close(sem)
	}()
	
	return results
}

func Wait[T any](results <-chan *Result[T], handler func(v T) error) error {
	for item := range results {
		if item.Error != nil {
			return item.Error
		}
		if handler != nil {
			err := handler(item.Value)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}
