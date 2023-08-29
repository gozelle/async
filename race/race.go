package race

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gozelle/async"
	"github.com/gozelle/async/multierr"
)

type Null = async.Null

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

	if ctx == nil {
		ctx = context.Background()
	}

	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()

	vr := async.NewValue[T]()
	errs := multierr.Errors{}

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

			r, e := f.Runner(ctx)
			if e != nil {
				errs.AddError(e)
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

	err = errs.Error()

	return
}
