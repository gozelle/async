package pipe

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"reflect"
	"time"
)

// Run 按管道处理任务，可提前终止
func Run(ctx context.Context, initial any, runners ...async.NamedRunner) (ret any, err error) {
	if reflect.TypeOf(initial).Kind() != reflect.Pointer {
		err = fmt.Errorf("initial expect a pointer")
		return
	}
	if len(runners) == 0 {
		err = fmt.Errorf("no runner")
		return
	}
	for _, v := range runners {
		var exit bool
		ret, exit, err = v.Runner(ctx, initial)
		if err != nil {
			return
		}
		if exit {
			break
		}
		initial = ret
	}
	return
}

// RunWithTiming 按管道处理任务，并统计每项任务执行的时间
// NamedRunner.Name 非必填，只用作耗时展示区分
func RunWithTiming(ctx context.Context, initial any, runners ...async.NamedRunner) (ret any, intervals []async.NamedInterval, err error) {
	if reflect.TypeOf(initial).Kind() != reflect.Pointer {
		err = fmt.Errorf("initial expect a pointer")
		return
	}
	if len(runners) == 0 {
		err = fmt.Errorf("no runner")
		return
	}
	for _, v := range runners {
		var exit bool
		now := time.Now()
		ret, exit, err = v.Runner(ctx, initial)
		intervals = append(intervals, async.NamedInterval{Name: v.Name, Interval: time.Since(now)})
		if err != nil {
			return
		}
		if exit {
			break
		}
		initial = ret
	}
	return
}
