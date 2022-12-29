package race

import (
	"context"
	"fmt"
	"github.com/gozelle/godash"
	"sync"
	"time"
)

func Race(ctx context.Context, handlers ...godash.Handler) (result any, err error) {
	
	if len(handlers) == 0 {
		err = fmt.Errorf("no handlers")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := godash.NewValue()
	ve := godash.NewValue()
	
	wg := sync.WaitGroup{}
	wg.Add(len(handlers))
	
	for _, f := range handlers {
		go func(f godash.Handler) {
			go func() {
				select {
				case <-cctx.Done():
					wg.Done()
				}
			}()
			r, e := f(cctx)
			if e != nil {
				if ve.GetValue() == nil {
					ve.SetValue(err)
				}
				return
			}
			if vr.GetValue() == nil {
				vr.SetValue(r)
			}
			cancel()
		}(f)
	}
	wg.Wait()
	
	result = vr.GetValue()
	if e := ve.GetValue(); e != nil {
		err = e.(error)
	}
	
	return
}

func DelayRace(ctx context.Context, handlers ...*godash.DelayHandler) (result any, err error) {
	
	if len(handlers) == 0 {
		err = fmt.Errorf("no handlers")
		return
	}
	
	cctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	
	vr := godash.NewValue()
	ve := godash.NewValue()
	
	wg := sync.WaitGroup{}
	wg.Add(len(handlers))
	
	for _, f := range handlers {
		go func(f *godash.DelayHandler) {
			go func() {
				select {
				case <-cctx.Done():
					wg.Done()
				}
			}()
			time.Sleep(f.Delay)
			if ctx.Err() != nil {
				return
			}
			r, e := f.Handler(cctx)
			if e != nil {
				if ve.GetValue() == nil {
					ve.SetValue(err)
				}
				return
			}
			if vr.GetValue() == nil {
				vr.SetValue(r)
			}
			cancel()
		}(f)
	}
	wg.Wait()
	
	result = vr.GetValue()
	if e := ve.GetValue(); e != nil {
		err = e.(error)
	}
	
	return
}
