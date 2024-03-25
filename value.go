package async

import (
	"sync"
)

func NewValue[T any]() *Value[T] {
	return &Value[T]{}
}

type Value[T any] struct {
	value    T
	nonempty bool
	lock     sync.RWMutex
}

func (r *Value[T]) Empty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return !r.nonempty
}

func (r *Value[T]) SetValue(val T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.value = val
	r.nonempty = true
}

func (r *Value[T]) SetValueOnce(val T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.nonempty {
		return
	}
	r.value = val
	r.nonempty = true
}

func (r *Value[T]) Value() T {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.value
}

func NewValues[T any]() *Values[T] {
	return &Values[T]{}
}

type Values[T any] struct {
	lock   sync.Mutex
	values []T
}

func (r *Values[T]) AddValue(val T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.values = append(r.values, val)
}

func (r *Values[T]) Values() []T {
	r.lock.Lock()
	defer r.lock.Unlock()

	var values []T
	for _, v := range r.values {
		values = append(values, v)
	}

	return values
}

func (r *Values[T]) Empty() bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	return len(r.values) == 0
}
