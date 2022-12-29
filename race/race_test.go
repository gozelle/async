package race

import (
	"context"
	"github.com/gozelle/flow"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestRace(t *testing.T) {
	r, err := Race(context.Background(), func(ctx context.Context) (result any, err error) {
		time.Sleep(500 * time.Millisecond)
		result = 1
		return
	}, func(ctx context.Context) (result any, err error) {
		time.Sleep(2000 * time.Millisecond)
		result = 2
		return
	})
	require.NoError(t, err)
	require.Equal(t, r.(int), 1)
}

func TestDelayRace(t *testing.T) {
	r, err := DelayRace(context.Background(), &godash.DelayHandler{
		Delay: 500 * time.Millisecond,
		Handler: func(ctx context.Context) (result any, err error) {
			result = 1
			return
		},
	}, &godash.DelayHandler{
		Delay: 2000 * time.Millisecond,
		Handler: func(ctx context.Context) (result any, err error) {
			result = 2
			return
		},
	})
	require.NoError(t, err)
	require.Equal(t, r.(int), 1)
}
