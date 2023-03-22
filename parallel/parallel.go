package parallel

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/async/chunk"
	"sync"
)

type Runner = async.Runner

// Run 并发执行 Runners
// 效果：Runners 产生的结果将无序存放在 result 中，如果有 Runner 产生了错误，err 会存放第 1 个产生的 error
func Run(ctx context.Context, runners ...Runner) (result []any, err error) {
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	l := len(runners)
	ev := async.NewValue()
	rv := async.NewValues()
	
	wg := sync.WaitGroup{}
	wg.Add(l)
	
	for _, v := range runners {
		go func(runner async.Runner) {
			defer func() {
				wg.Done()
			}()
			r, e := runner(ctx)
			if e != nil {
				if ev.GetValue() == nil {
					ev.SetValue(e)
				}
				return
			}
			rv.AddValue(r)
		}(v)
	}
	wg.Wait()
	
	if ev.GetValue() != nil {
		err = ev.GetValue().(error)
		return
	}
	
	result = rv.GetValues()
	
	return
}

// RunWithChunk 分片并发任务
// 效果：将会按照 chunks 定义的数值，使用回调函数对任务进行分组处理
// 回调函数会同步运行,如果回调函数返回错误，将会终止 RunWithChunk 运行, 返回错误；
// 如果希望在分组全部运行结束后，再提交回调函数处理后的结果，请在外面做好事务控制。
func RunWithChunk(ctx context.Context, chunks int, callback func(values []any) error, runners ...Runner) (err error) {
	
	if chunks <= 0 {
		err = fmt.Errorf("chunks expect greater than 0")
		return
	}
	
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	ranges, err := chunk.SplitInt64s(0, int64(len(runners))-1, int64(chunks))
	if err != nil {
		err = fmt.Errorf("calc ranges error: %s", err)
		return
	}
	
	for _, r := range ranges {
		var execRunners []async.Runner
		for i := r.Begin; i <= r.End; i++ {
			execRunners = append(execRunners, runners[int(i)])
		}
		var values []any
		values, err = Run(ctx, execRunners...)
		if err != nil {
			return
		}
		if callback != nil {
			err = callback(values)
			if err != nil {
				err = fmt.Errorf("exec callback error: %s", err)
				return
			}
		}
	}
	
	return
}
