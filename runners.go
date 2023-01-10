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

type NamedRunner struct {
	Name   string
	Runner func(ctx context.Context, initial any) (value any, exit bool, err error)
}
