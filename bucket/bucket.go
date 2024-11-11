package bucket

import (
	"fmt"
	"sync"
	"time"
)

type Handler[T any] func(done <-chan struct{}, data []T)

func NewBucket[T any](threshold uint, interval time.Duration, handler Handler[T]) *Bucket[T] {
	return &Bucket[T]{
		data:      make([]T, 0),
		threshold: threshold,
		interval:  interval,
		handler:   handler,
		done:      make(chan struct{}),
	}
}

type Bucket[T any] struct {
	lock      sync.Mutex
	data      []T           // 存放数据的 slice
	interval  time.Duration // 处理时间间隔
	handler   Handler[T]
	done      chan struct{}
	closed    bool
	threshold uint
}

func (b *Bucket[T]) Stop() {
	b.done <- struct{}{}
	b.closed = true
}

// Push 仅当桶处于 closed 状态时会报错
func (b *Bucket[T]) Push(data T) (err error) {

	if b.closed {
		err = fmt.Errorf("bucket has closed")
		return
	}

	b.lock.Lock()
	defer func() {
		b.lock.Unlock()
	}()

	b.data = append(b.data, data)

	return
}

func (b *Bucket[T]) Start() {

	defer func() {
		close(b.done)
	}()

	b.closed = false
	timer := time.NewTimer(b.interval)
	defer func() {
		timer.Stop()
	}()

	for {
		select {
		case <-b.done:
			return
		case <-timer.C:
			if len(b.data) > 0 {
				b.process(timer)
			}
		default:
			if uint(len(b.data)) >= b.threshold {
				b.process(timer)
			}
		}
	}
}

func (b *Bucket[T]) process(timer *time.Timer) {

	b.lock.Lock()
	defer func() {
		b.lock.Unlock()
	}()

	timer.Stop()
	if b.handler != nil {
		b.handler(b.done, append([]T{}, b.data...))
	}
	b.data = make([]T, 0)
	timer.Reset(b.interval)
}
