package parallel

import (
	"context"
	"fmt"
	"github.com/gozelle/async"
	"github.com/gozelle/testify/require"
	"testing"
)

func TestRun(t *testing.T) {
	
	var runners []async.Runner
	
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			result = v
			return
		})
	}
	
	values, err := Run(context.Background(), runners...)
	require.NoError(t, err)
	n := 0
	for _, v := range values {
		n += v.(int)
	}
	require.Equal(t, 55, n)
}

func TestRunLimit(t *testing.T) {
	var runners []async.Runner
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			result = v
			return
		})
	}
	n := 0
	err := RunWithChunk(context.Background(), 3, func(values []any) error {
		for _, v := range values {
			n += v.(int)
		}
		return nil
	}, runners...)
	require.NoError(t, err)
	require.Equal(t, 55, n)
}

func TestRunError(t *testing.T) {
	
	var runners []async.Runner
	
	for i := 1; i <= 5; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			if v == 3 {
				err = fmt.Errorf("some error")
				return
			}
			result = v
			return
		})
	}
	
	_, err := Run(context.Background(), runners...)
	require.Error(t, err)
}

func TestRunLimitError(t *testing.T) {
	var runners []async.Runner
	for i := 1; i <= 10; i++ {
		v := i
		runners = append(runners, func(ctx context.Context) (result any, err error) {
			if v == 7 {
				err = fmt.Errorf("some error")
				return
			}
			result = v
			return
		})
	}
	n := 0
	err := RunWithChunk(context.Background(), 3, func(values []any) error {
		for _, v := range values {
			n += v.(int)
		}
		return nil
	}, runners...)
	require.Error(t, err)
}
