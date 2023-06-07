package bucket_test

import (
	"testing"
	"time"

	"github.com/gozelle/async/bucket"
)

func TestBatch(t *testing.T) {

	b := bucket.NewBucket[int](2*time.Second, 500, func(done <-chan struct{}, data []int) {
		t.Log(len(data))
	})

	go b.Start()

	go func() {

	}()

	i := 0
	for {
		i++
		if i > 100000000 {
			break
		}
		b.Push(i)
		//time.Sleep(1 * time.Millisecond)
		//if i%100 == 0 {
		//	time.Sleep(2 * time.Second)
		//}
	}
	select {}
}
