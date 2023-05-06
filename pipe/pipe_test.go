package pipe

import (
	"context"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestPipe(t *testing.T) {
	var r Runner[int] = func(ctx context.Context, initial *int) error {
		*initial += 1
		return nil
	}
	a := 1
	err := Run(context.Background(), &a, []Runner[int]{r})
	require.NoError(t, err)
	require.Equal(t, 2, a)
}
