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

type IRunner[T any] interface {
	Runner() Runner[T]
}

type Runner[T any] struct {
	Delay  time.Duration
	Runner async.Runner[T]
}

// Run
// 并发执行 Runner, 返回其中最快的结果, 忽略其它较慢的结果或错误
// 如果全部返回错误，则返回出现的第一个错误
// 配置延迟执行的 Runner，会在等待配置时间后，再开始执行
func Run[T any](ctx context.Context, runners []*Runner[T]) (result T, err error) {
	
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := async.NewValue[T]()
	ve := async.NewValues[error]()
	
	wg := sync.WaitGroup{}
	
	for _, f := range runners {
		wg.Add(1)
		go func(f *Runner[T]) {
			
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
			
			if f.Delay > 0 {
				time.Sleep(f.Delay)
			}
			
			if cctx.Err() != nil {
				return
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
	if !vr.Empty() {
		result = vr.Value()
		return
	}
	
	err = multierr.Combine(ve.Values()...)
	
	return
}
