package pipe

import (
	"context"
	"fmt"
	"reflect"
)

type IRunner[T any] interface {
	Runner() Runner[T]
}

type Runner[T any] func(ctx context.Context, initial *T) (err error)

// Run 按管道处理任务，可提前终止
func Run[T any](ctx context.Context, initial *T, runners []Runner[T]) (err error) {
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
