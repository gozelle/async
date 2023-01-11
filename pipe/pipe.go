package pipe

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/logging"
	"reflect"
	"time"
)

var log = logging.Logger("pipe")

// Run 按管道处理任务，可提前终止
func Run(ctx context.Context, initial any, runners ...async.NamedPipeRunner) (err error) {
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
		exit, err = v.Runner(ctx, initial)
		if err != nil {
			return
		}
		if exit {
			break
		}
	}
	return
}

// RunWithTiming 按管道处理任务，并统计每项任务执行的时间
// NamedRunner.Name 非必填，只用作耗时展示区分
func RunWithTiming(ctx context.Context, initial any, runners ...async.NamedPipeRunner) (elapseds []async.NamedElapsed, err error) {
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
		exit, err = v.Runner(ctx, initial)
		elapsed := time.Since(now)
		log.Debugf("exec '%s' elapsed time: %s", v.Name, elapsed)
		elapseds = append(elapseds, async.NamedElapsed{Name: v.Name, Elapsed: elapsed})
		if err != nil {
			return
		}
		if exit {
			break
		}
	}
	return
}
