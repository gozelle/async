package race

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/multierr"
	"sync"
	"sync/atomic"
	"time"
)

type IRunner interface {
	Runner() Runner
}

type Runner struct {
	Delay  time.Duration
	Runner async.Runner
}

// Run
// 并发执行 Runner, 返回其中最快的结果
// 如果全部返回错误，则返回出现的第一个错误
// 配置延迟执行的 Runner，会在等待配置时间后，再开始执行
func Run[T any](ctx context.Context, runners []*Runner) (result T, err error) {
	
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := async.NewValue()
	ve := async.NewValues()
	
	wg := sync.WaitGroup{}
	
	for _, f := range runners {
		wg.Add(1)
		go func(f *Runner) {
			
			done := atomic.Value{}
			go func() {
				select {
				case <-cctx.Done():
					
					if done.Load() == nil {
						wg.Done()
						done.Store(true)
					}
				}
			}()
			
			defer func() {
				if done.Load() == nil {
					wg.Done()
					done.Store(true)
				}
			}()
			
			if cctx.Err() != nil {
				return
			}
			if f.Delay > 0 {
				time.Sleep(f.Delay)
			}
			
			r, e := f.Runner(cctx)
			if e != nil {
				ve.AddValue(e)
				return
			}
			vr.SetValue(r)
			cancel()
			
			return
		}(f)
	}
	wg.Wait()
	
	if !ve.Empty() {
		list := ve.GetValues()
		errors := make([]error, len(list))
		for _, e := range list {
			errors = append(errors, e.(error))
		}
		err = multierr.Combine(errors...)
		return
	}
	
	vv, ok := vr.GetValue().(T)
	if !ok {
		err = fmt.Errorf("can't assert value: %v to type: T", vr.GetValue())
		return
	}
	result = vv
	
	return
}
