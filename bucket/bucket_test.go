package bucket_test

import (
	"testing"
	"time"
)

func TestBatch(t *testing.T) {

	b := bucket.NewBatch[int](2*time.Second, 500, func(done <-chan struct{}, data []int) {
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
		b.Add(i)
		//time.Sleep(1 * time.Millisecond)
		//if i%100 == 0 {
		//	time.Sleep(2 * time.Second)
		//}
	}
	select {}
}
