package async

import (
	"context"
)

type IRunner[T any] interface {
	Runner() Runner[T]
}

type Runner[T any] func(ctx context.Context) (T, error)
