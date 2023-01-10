package forever

import (
	"context"
	"fmt"
	"time"
)

// Run 执行永久任务，类似 For 循环
// 当 runner 返回时 error 时，将会退出循环
func Run(ctx context.Context, interval time.Duration, runner func(ctx context.Context) error) (err error) {
	if interval < time.Second {
		err = fmt.Errorf("interval expect greater than 1 second")
		return
	}
	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()
	for {
		select {
		case <-timer.C:
			err = runner(ctx)
			if err != nil {
				return
			}
			timer.Reset(interval)
		case <-ctx.Done():
			return
		}
	}
}
