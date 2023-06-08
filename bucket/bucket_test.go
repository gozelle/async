package bucket_test

import (
	"log"
	"testing"
	"time"
	
	"github.com/gozelle/async/bucket"
)

func TestBatch(t *testing.T) {
	
	b := bucket.NewBucket[int](500, 200*time.Millisecond, func(done <-chan struct{}, data []int) {
		log.Printf("%d", len(data))
	})
	go func() {
		i := 0
		for {
			i++
			if i > 1000000 {
				b.Stop()
				break
			}
			err := b.Push(i)
			if err != nil {
				t.Log(err)
				break
			}
		}
	}()
	b.Start()
}
