package semaphore

import (
	"context"
	"golang.org/x/sync/semaphore"
)

func WithContext(ctx context.Context, limit int64) *Semaphore {
	return &Semaphore{
		Weighted: semaphore.NewWeighted(limit),
	}
}

type Semaphore struct {
	*semaphore.Weighted
}

func (s *Semaphore) Acquire(ctx context.Context, n int64) error {
	return s.Weighted.Acquire(ctx, n)
}

func (s *Semaphore) TryAcquire(n int64) bool {
	return s.Weighted.TryAcquire(n)
}

func (s *Semaphore) Release(n int64) {
	s.Weighted.Release(n)
}
