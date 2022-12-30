package race

import (
	"context"
	"github.com/gozelle/godash"
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
	handlers := []*godash.DelayHandler{
		{
			Delay: 0,
			Handler: func(ctx context.Context) (result any, err error) {
				result = 1
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Handler: func(ctx context.Context) (result any, err error) {
				result = 3
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Handler: func(ctx context.Context) (result any, err error) {
				result = 2
				return
			},
		},
	}
	
	var handlers2 []*godash.DelayHandler
	for _, h := range handlers {
		v := h
		func(v *godash.DelayHandler) {
			handlers2 = append(
				handlers2,
				&godash.DelayHandler{
					Delay: v.Delay,
					Handler: func(ctx context.Context) (result any, err error) {
						return v.Handler(ctx)
					},
				},
			)
		}(v)
	}
	
	r, err := DelayRace(context.Background(), handlers2...)
	require.NoError(t, err)
	require.Equal(t, r.(int), 1)
}
