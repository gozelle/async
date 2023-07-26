package parallel

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	
	g := NewGroup()
	g.SetLimit(1)
	g.Go(func() error {
		return fmt.Errorf("some error")
	})
	var a int
	g.Go(func() error {
		a = 1
		return nil
	})
	err := g.Wait()
	t.Log(err)
	t.Log(a)
}
