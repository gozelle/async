package retry

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"time"
)

func Run(ctx context.Context, times int, interval time.Duration, runner async.Runner) (result any, err error) {
	if times < 1 {
		err = fmt.Errorf("times expact > 1,got: %d", times)
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
