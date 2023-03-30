package race_test

import (
	"context"
	"fmt"
	"github.com/gozelle/async/race"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestDelayRace(t *testing.T) {
	runners := []*race.Runner[int]{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 1
				t.Log(result)
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				t.Log(result)
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				t.Log(result)
				return
			},
		},
	}
	
	r, err := race.Run[int](context.Background(), runners)
	require.NoError(t, err)
	require.Equal(t, 1, r)
}

func TestRaceError(t *testing.T) {
	runners := []*race.Runner[int]{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 1
				err = fmt.Errorf("some error form: 1")
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				err = fmt.Errorf("some error form: 3")
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				err = fmt.Errorf("some error form: 2")
				return
			},
		},
	}
	
	_, err := race.Run[int](context.Background(), runners)
	require.Error(t, err)
	t.Log(err)
}
