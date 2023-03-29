package race

import (
	"context"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestDelayRace(t *testing.T) {
	handlers := []*Runner{
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
	
	var handlers2 []*Runner
	for _, h := range handlers {
		v := h
		func(v *Runner) {
			handlers2 = append(
				handlers2,
				&Runner{
					Delay: v.Delay,
					Runner: func(ctx context.Context) (result any, err error) {
						return v.Runner(ctx)
					},
				},
			)
		}(v)
	}
	
	r, err := Run(context.Background(), handlers2)
	require.NoError(t, err)
	require.Equal(t, r.(int), 1)
}
