package race

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"sync"
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
func Run(ctx context.Context, runners []*Runner) (result any, err error) {
	
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := async.NewValue()
	ve := async.NewValue()
	
	wg := sync.WaitGroup{}
	wg.Add(len(runners))
	
	for _, f := range runners {
		go func(f *Runner) {
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
