package race

import (
	"context"
	"github.com/gozelle/async"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestRace(t *testing.T) {
	r, err := Run(context.Background(), func(ctx context.Context) (result any, err error) {
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
	handlers := []*async.DelayRunner{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 1
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 3
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result any, err error) {
				result = 2
				return
			},
		},
	}
	
	var handlers2 []*async.DelayRunner
	for _, h := range handlers {
		v := h
		func(v *async.DelayRunner) {
			handlers2 = append(
				handlers2,
				&async.DelayRunner{
					Delay: v.Delay,
					Runner: func(ctx context.Context) (result any, err error) {
						return v.Runner(ctx)
					},
				},
			)
		}(v)
	}
	
	r, err := RunWithDelay(context.Background(), handlers2...)
	require.NoError(t, err)
	require.Equal(t, r.(int), 1)
}
