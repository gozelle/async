package forever

import (
	"context"
	"fmt"
	"time"
)

// Run 执行永久任务，类似 For 循环
// 当 runner 返回时 error 时，将会退出循环
func Run(ctx context.Context, duration time.Duration, runner func(ctx context.Context) error) (err error) {
	if duration < time.Second {
		err = fmt.Errorf("duration expect greater than 1 second")
		return
	}
	ticker := time.NewTicker(duration)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			err = runner(ctx)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
