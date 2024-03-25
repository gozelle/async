package multierr

import (
	"fmt"
	"sync"
)

type Errors struct {
	lock   sync.Mutex
	errors map[error]struct{}
}

func (e *Errors) AddError(err error) {
	e.lock.Lock()
	defer func() {
		e.lock.Unlock()
	}()
	if err == nil {
		return
	}
	if e.errors == nil {
		e.errors = make(map[error]struct{})
	}
	e.errors[err] = struct{}{}
}

func (e *Errors) Error() error {
	e.lock.Lock()
	defer func() {
		e.lock.Unlock()
	}()
	var err error
	if e.errors != nil {
		for v := range e.errors {
			if err == nil {
				err = v
			} else {
				err = fmt.Errorf("%v; %w", err, v)
			}
		}
	}
	return err
}

func (e *Errors) String() string {
	err := e.Error()
	if err == nil {
		return ""
	}
	return err.Error()
}
