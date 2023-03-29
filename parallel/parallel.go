package parallel

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"sync"
)

type IRunner interface {
	Runner() Runner
}

type Runner = async.Runner

func Run(ctx context.Context, limit uint, runners []Runner) <-chan any {
	
	results := make(chan any, len(runners))
	
	if limit == 0 {
		defer func() {
			results <- fmt.Errorf("limit expect great than 0")
			close(results)
		}()
		return results
	}
	
	cctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	
	sem := make(chan struct{}, limit)
	
	for _, v := range runners {
		wg.Add(1)
		sem <- struct{}{}
		go func(runner Runner) {
			defer func() {
				<-sem
				wg.Done()
			}()
			select {
			case <-cctx.Done():
				return
			default:
				result, err := runner(ctx)
				if err != nil {
					results <- err
					cancel()
				} else {
					results <- result
				}
			}
		}(v)
	}
	
	go func() {
		wg.Wait()
		close(results)
		cancel()
	}()
	
	return results
}
