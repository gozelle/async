package async

import (
	"context"
)

type IRunner interface {
	Runner() Runner
}

type Runner func(ctx context.Context) (result any, err error)
