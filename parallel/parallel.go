package parallel

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/async/chunk"
	"sync"
)

// Run 并发执行 Runners
// 效果：Runners 产生的结果将无序存放在 result 中，如果有 Runner 产生了错误，err 会存放第 1 个产生的 error
// TODO: 加入有错误产生时，提前退出
func Run(ctx context.Context, runners ...async.Runner) (result []any, err error) {
	
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

// RunLimit 分片并发任务
// 效果：将会按照 limit 定义的数值，使用回调函数对任务进行分组处理
// 回调函数会同步运行,如果回调函数返回错误，将会终止 RunLimit 运行, 返回错误；
// 如果希望在分组全部运行结束后，再提交回调函数处理后的结果，请在外面做好事务控制。
func RunLimit(ctx context.Context, limit int, callback func(values []any) error, runners ...async.Runner) (err error) {
	
	if limit <= 0 {
		err = fmt.Errorf("limit expect greater than 0")
		return
	}
	
	if len(runners) == 0 {
		err = fmt.Errorf("no runners")
		return
	}
	
	ranges, err := chunk.Int64s(0, int64(len(runners))-1, int64(limit))
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
		err = callback(values)
		if err != nil {
			err = fmt.Errorf("exec callback error: %s", err)
			return
		}
	}
	
	return
}
