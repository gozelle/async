package forever

import (
	"context"
	"fmt"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestForever(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	err := Run(ctx, time.Second, func(ctx context.Context) error {
		t.Log(time.Now())
		return nil
	})
	require.NoError(t, err)
}

func TestForeverError(t *testing.T) {
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	err := Run(ctx, time.Second, func(ctx context.Context) error {
		t.Log(time.Now())
		if time.Since(now) > 2*time.Second {
			return fmt.Errorf("timeout")
		}
		return nil
	})
	require.Error(t, err)
}
