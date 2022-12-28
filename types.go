package flow

import (
	"context"
	"sync"
	"time"
)

type Handler func(ctx context.Context) (result any, err error)

type DelayHandler struct {
	Delay   time.Duration
	Handler Handler
}

func NewValue() *Value {
	return &Value{}
}

type Value struct {
	value any
	lock  sync.RWMutex
}

func (r *Value) SetValue(val any) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.value = val
}

func (r *Value) GetValue() any {
	r.lock.RLock()
	defer r.lock.RUnlock()
	
	return r.value
}
