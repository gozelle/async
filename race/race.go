package race

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"sync"
	"time"
)

// Run
// 并发执行 Runner, 返回其中最快的结果
// 如果全部返回错误，则返回出现的第一个错误
func Run(ctx context.Context, handlers ...async.Runner) (result any, err error) {
	
	if len(handlers) == 0 {
		err = fmt.Errorf("no handlers")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := async.NewValue()
	ve := async.NewValue()
	
	wg := sync.WaitGroup{}
	wg.Add(len(handlers))
	
	for _, f := range handlers {
		go func(f async.Runner) {
			go func() {
				select {
				case <-cctx.Done():
					wg.Done()
				}
			}()
			r, e := f(cctx)
			if e != nil {
				if ve.GetValue() == nil {
					ve.SetValue(err)
				}
				return
			}
			if vr.GetValue() == nil {
				vr.SetValue(r)
			}
			cancel()
		}(f)
	}
	wg.Wait()
	
	result = vr.GetValue()
	if e := ve.GetValue(); e != nil {
		err = e.(error)
	}
	
	return
}

// RunWithDelay
// 并发执行 Runner, 返回其中最快的结果
// 如果全部返回错误，则返回出现的第一个错误
// 与 Run 不同的是，配置延迟执行的 Runner，会在等待配置时间后，再开始执行
func RunWithDelay(ctx context.Context, handlers ...*async.DelayRunner) (result any, err error) {
	
	if len(handlers) == 0 {
		err = fmt.Errorf("no handlers")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := async.NewValue()
	ve := async.NewValue()
	
	wg := sync.WaitGroup{}
	wg.Add(len(handlers))
	
	for _, f := range handlers {
		go func(f *async.DelayRunner) {
			go func() {
				select {
				case <-cctx.Done():
					wg.Done()
				}
			}()
			time.Sleep(f.Delay)
			if ctx.Err() != nil {
				return
			}
			r, e := f.Runner(cctx)
			if e != nil {
				if ve.GetValue() == nil {
					ve.SetValue(err)
				}
				return
			}
			if vr.GetValue() == nil {
				vr.SetValue(r)
			}
			cancel()
		}(f)
	}
	wg.Wait()
	
	result = vr.GetValue()
	if e := ve.GetValue(); e != nil {
		err = e.(error)
	}
	
	return
}
