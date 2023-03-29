package pipe

import (
	"context"
	"fmt"
	"reflect"
)

type IRunner interface {
	Runner() Runner
}

type Runner func(ctx context.Context, initial any) (err error)

// Run 按管道处理任务，可提前终止
func Run(ctx context.Context, initial any, runners []Runner) (err error) {
	if reflect.TypeOf(initial).Kind() != reflect.Pointer {
		err = fmt.Errorf("initial expect a pointer")
		return
	}
	if len(runners) == 0 {
		err = fmt.Errorf("no runner")
		return
	}
	
	for _, v := range runners {
		err = v(ctx, initial)
		if err != nil {
			return
		}
	}
	return
}
