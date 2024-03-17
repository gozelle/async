package parallel

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/gozelle/async"
	"github.com/gozelle/async/multierr"
)

type Null = async.Null

type Result[T any] struct {
	Error error
	Value T
}

type Runner[T any] async.Runner[T]

func Run[T any](ctx context.Context, limit uint, runners []Runner[T]) <-chan *Result[T] {

	ch := make(chan *Result[T], len(runners))

	if limit == 0 {
		defer func() {
			ch <- &Result[T]{Error: fmt.Errorf("limit expect great than 0")}
			close(ch)
		}()
		return ch
	}

	if ctx == nil {
		ctx = context.Background()
	}

	go run[T](ctx, limit, runners, ch)

	return ch
}
func run[T any](ctx context.Context, limit uint, runners []Runner[T], ch chan *Result[T]) {

	errs := multierr.Errors{}
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, limit)

	defer func() {
		close(ch)
		close(sem)
	}()

	for _, v := range runners {

		// achieve a blocking effect by sending semaphores to a channel with a specified capacity of "limit"
		// when the channel is full, it will block here until a task is completed and frees up channel capacity
		sem <- struct{}{}

		// if the semaphore is acquired, prioritize checking whether the context has done.
		// if it has, break out of the for loop.
		select {
		case <-ctx.Done():
			errs.AddError(ctx.Err())
			<-sem
		default:
			// when an error occurs, the semaphores of all subsequent tasks will be directly ignored.
			if errs.Error() != nil {
				<-sem
				continue
			}
			wg.Add(1)
			go func(runner Runner[T]) {
				defer func() {
					e := recover()
					if e != nil {
						errs.AddError(fmt.Errorf("%v", e))
						debug.PrintStack()
					}
					// the task has been executed to completion,
					// release the semaphore.
					<-sem
					wg.Done()
				}()

				r, e := runner(ctx)
				if e != nil {
					errs.AddError(e)
				} else {
					ch <- &Result[T]{Value: r}
				}
			}(v)
		}
	}

	wg.Wait()

	// all tasks have been completed.
	// check for any errors and ensure that the error is the last result sent to the channel.
	if errs.Error() != nil {
		ch <- &Result[T]{Error: errs.Error()}
	}
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
