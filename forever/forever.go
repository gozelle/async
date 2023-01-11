package forever

import (
	"time"
)

// Run 执行永久任务，类似 For 循环
// 当 runner 返回时 error 时，将会退出循环
func Run(interval time.Duration, runner func()) {
	if interval < time.Second {
		interval = time.Second
	}
	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()
	for {
		select {
		case <-timer.C:
			runner()
			timer.Reset(interval)
		}
	}
}
