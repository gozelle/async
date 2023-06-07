package batch

import (
	"fmt"
	"sync"
	"time"
)

type Handler[T any] func(done <-chan struct{}, data []T)

func NewBatch[T any](interval time.Duration, threshold uint, handler Handler[T]) *Batch[T] {
	return &Batch[T]{
		data:      make([]T, 0),
		interval:  interval,
		handler:   handler,
		done:      make(chan struct{}),
		threshold: threshold,
	}
}

type Batch[T any] struct {
	lock      sync.Mutex
	data      []T           // 存放数据的 slice
	interval  time.Duration // 处理时间间隔
	handler   Handler[T]
	done      chan struct{}
	closed    bool
	threshold uint
}

func (b *Batch[T]) Stop() {
	b.done <- struct{}{}
	b.closed = true
}

func (b *Batch[T]) Add(data T) (err error) {

	if b.closed {
		err = fmt.Errorf("batch has closed")
		return
	}

	b.lock.Lock()
	defer func() {
		b.lock.Unlock()
	}()

	b.data = append(b.data, data)

	return
}

func (b *Batch[T]) Start() {

	defer func() {
		close(b.done)
	}()

	b.closed = false
	ticker := time.NewTicker(b.interval)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-b.done:
			break
		case <-ticker.C:
			if len(b.data) > 0 {
				b.process()
			}
		default:
			if uint(len(b.data)) >= b.threshold {
				b.process()
			}
		}
	}
}

func (b *Batch[T]) process() {

	b.lock.Lock()
	defer func() {
		b.lock.Unlock()
	}()

	if b.handler != nil {
		b.handler(b.done, append([]T{}, b.data...))
	}

	b.data = nil
}
