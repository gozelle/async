package async

import (
	"context"
	"time"
)

type Runner func(ctx context.Context) (result any, err error)

type DelayRunner struct {
	Delay  time.Duration
	Runner Runner
}

// Int64Range 记录 Int64 区间
type Int64Range struct {
	Begin int64
	End   int64
}

// Len 返回 Range 的有效长度
// Begin=0, End=0  则： Len = 1
// Begin=1, End=2  则:  Len = 2
func (i Int64Range) Len() int64 {
	return i.End - i.Begin + 1
}
