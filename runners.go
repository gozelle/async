package async

import (
	"context"
	"time"
)

type IRunner interface {
	Runner() Runner
}

type IDelayRunner interface {
	DelayRunner() DelayRunner
}

type IPipeRunner interface {
	PipeRunner() PipeRunner
}

type INamedPipeRunner interface {
	NamedPipeRunner() NamedPipeRunner
}

type Runner func(ctx context.Context) (result any, err error)

type DelayRunner struct {
	Delay  time.Duration
	Runner Runner
}

type NamedRunner struct {
	Name   string
	Runner func(ctx context.Context) (value any, err error)
}

type PipeRunner func(ctx context.Context, initial any) (err error)

type NamedPipeRunner struct {
	Name   string
	Runner func(ctx context.Context, initial any) (err error)
}
