package retry

import (
	"context"
	"fmt"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	now := time.Now()
	r, err := Run(context.Background(), 3, 10*time.Second, func(ctx context.Context) (result any, err error) {
		t.Log(time.Now())
		d := time.Since(now)
		if d < 2*time.Second {
			err = fmt.Errorf("some error")
		} else {
			result = d
		}
		return
	})
	require.NoError(t, err)
	t.Log(r)
}
