package forever

import (
	"testing"
	"time"
)

func TestForever(t *testing.T) {
	Run(time.Second, func() {
		t.Log(time.Now())
	})
}
