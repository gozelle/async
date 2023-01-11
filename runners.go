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
	Runner func(ctx context.Context) (value any, err error)
}

type PipeRunner func(ctx context.Context, initial any) (exit bool, err error)

type NamedPipeRunner struct {
	Name   string
	Runner func(ctx context.Context, initial any) (exit bool, err error)
}
