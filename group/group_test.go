package group

import (
	"fmt"
	"github.com/gozelle/testify/require"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	g := NewGroup()
	g.SetLimit(1)
	g.Go(func() error {
		return fmt.Errorf("error 2")
	})
	var a int
	g.Go(func() error {
		time.Sleep(3 * time.Second)
		a = 1
		return nil
	})
	err := g.Wait()
	t.Log(err)
	require.Error(t, err)
	require.Equal(t, 0, a)
}
