package race_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gozelle/async/race"
	"github.com/gozelle/testify/require"
)

func TestDelayRace(t *testing.T) {
	runners := []*race.Runner[int]{
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 0
				time.Sleep(5 * time.Second)
				return
			},
		},
		{
			Delay: 0,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 0
				err = fmt.Errorf("some error")
				return
			},
		},
		{
			Delay: 1,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 1
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				return
			},
		},
	}

	for i := 0; i < 100000; i++ {
		r, err := race.Run[int](context.Background(), runners)
		require.NoError(t, err)
		require.Equal(t, 1, r)
	}
}

func TestRaceError(t *testing.T) {
	runners := []*race.Runner[int]{
		{
			Delay: 1,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 1
				err = fmt.Errorf("some error form: 1")
				return
			},
		},
		{
			Delay: 2 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				err = fmt.Errorf("some error form: 2")
				return
			},
		},
		{
			Delay: 3 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				err = fmt.Errorf("some error form: 3")
				return
			},
		},
	}

	_, err := race.Run[int](context.Background(), runners)
	require.Error(t, err)

	require.True(t, strings.Contains(err.Error(), "some error form: 1"))
	require.True(t, strings.Contains(err.Error(), "some error form: 2"))
	require.True(t, strings.Contains(err.Error(), "some error form: 3"))
}

func TestContextCancel(t *testing.T) {
	runners := []*race.Runner[int]{
		{
			Runner: func(ctx context.Context) (result int, err error) {
				time.Sleep(10 * time.Second)
				result = 1
				return
			},
		},
		{
			Delay: 20 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 2
				err = fmt.Errorf("some error form: 2")
				return
			},
		},
		{
			Delay: 30 * time.Second,
			Runner: func(ctx context.Context) (result int, err error) {
				result = 3
				return
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		cancel()
	}()
	_, err := race.Run[int](ctx, runners)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
