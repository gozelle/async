package multierr

import (
	"fmt"
	"sync"
)

type Errors struct {
	lock   sync.Mutex
	errors []error
}

func (e *Errors) AddError(err error) {
	e.lock.Lock()
	defer func() {
		e.lock.Unlock()
	}()
	if err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e *Errors) Error() error {
	e.lock.Lock()
	defer func() {
		e.lock.Unlock()
	}()
	var err error
	if e.errors != nil {
		for _, v := range e.errors {
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
