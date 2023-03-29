package async

import (
	"sync"
	"time"
)

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

func NewValues() *Values {
	return &Values{}
}

type Values struct {
	lock   sync.Mutex
	values []any
}

func (r *Values) AddValue(val any) {
	r.lock.Lock()
	defer r.lock.Unlock()
	
	r.values = append(r.values, val)
}

func (r *Values) GetValues() []any {
	r.lock.Lock()
	defer r.lock.Unlock()
	
	var values []any
	for _, v := range r.values {
		values = append(values, v)
	}
	return values
}

func (r *Values) Empty() bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	
	return len(r.values) == 0
}

type NamedElapsed struct {
	Name    string
	Elapsed time.Duration
}
