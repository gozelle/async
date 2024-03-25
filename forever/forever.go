package forever

import (
	"context"
	"time"
)

// Run 执行永久任务，类似 For 循环
// 当 runner 返回时 error 时，将会退出循环
func Run(ctx context.Context, interval time.Duration, runner func()) {
	if interval < time.Second {
		interval = time.Second
	}
	runner()
	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			runner()
			timer.Reset(interval)
		}
	}
}
