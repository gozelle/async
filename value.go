package async

import (
	"sync"
)

func NewValue[T any]() *Value[T] {
	return &Value[T]{}
}

type Value[T any] struct {
	value T
	lock  sync.RWMutex
}

func (r *Value[T]) SetValue(val T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.value = val
}

func (r *Value[T]) GetValue() T {
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

func (r *Values[T]) GetValues() []T {
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
