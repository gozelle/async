package retry

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"time"
)

type IRunner[T any] interface {
	Runner() Runner[T]
}

type Runner[T any] async.Runner[T]

func Run[T any](ctx context.Context, times int, interval time.Duration, runner Runner[T]) (result T, err error) {
	if times < 1 {
		err = fmt.Errorf("times expact > 1, got: %d", times)
		return
	}
	for i := 0; i < times; i++ {
		result, err = runner(ctx)
		if err == nil {
			return
		}
		if interval > 0 {
			time.Sleep(interval)
		}
	}
	return
}
