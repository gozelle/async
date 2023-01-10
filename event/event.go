package event

import (
	"fmt"
	"log"
	"sync"
)

var (
	maxSeq    = 0
	listeners = make(map[string][]*listener) // map[事件id][]handleFunc
	lock      = sync.Mutex{}
)

type listener struct {
	seq     int
	handler func(args ...any)
}

// Subscribe 订阅时间
func Subscribe(id string, handler func(args ...any)) (seq int) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	
	maxSeq++
	seq = maxSeq
	ln := &listener{
		seq:     seq,
		handler: handler,
	}
	ls, ok := listeners[id]
	if !ok {
		ls = []*listener{}
	}
	listeners[id] = append(ls, ln)
	return seq
}

// SyncPublish 同步触发事件，若有处理器处理失败，则返回错误
func SyncPublish(id string, args ...any) (err error) {
	
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	
	ls, ok := listeners[id]
	if !ok {
		return
	}
	
	for _, l := range ls {
		e := catchPanic(func() {
			l.handler(args...)
		})
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
	}
	
	return
}

// AsyncPublish 异步触发事件，忽略处理器执行的结果
// 异步的处理器之间将会顺序同步执行
func AsyncPublish(id string, args ...any) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	
	ls, ok := listeners[id]
	if !ok {
		return
	}
	var rs []*listener
	for _, v := range ls {
		rs = append(rs, v)
	}
	go func(rs []*listener) {
		for _, r := range rs {
			_ = catchPanic(func() {
				r.handler(args...)
			})
		}
	}(rs)
	
	return
}

// Unsubscribe 取消订阅时间
func Unsubscribe(id string, seq int) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	
	ls, ok := listeners[id]
	if !ok {
		return
	}
	index := -1
	for i, l := range ls {
		if l.seq == seq {
			index = i
			break
		}
	}
	if index >= 0 {
		listeners[id] = append(ls[:index], ls[index+1:]...)
	}
}

func catchPanic(f func()) (err any) {
	defer func() {
		err = recover()
		if err != nil {
			log.Printf("[github.com/gozelle/event] CatchPanic panic: err=%v", err)
		}
	}()
	f()
	return
}
