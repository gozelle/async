package forever

import (
	"context"
	"testing"
	"time"

	"github.com/gozelle/testify/require"
)

func TestForever(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	count := 0
	Run(ctx, time.Second, func() {
		count++
	})
	require.Equal(t, 5, count)
}
