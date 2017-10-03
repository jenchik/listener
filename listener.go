package listener

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type (
	listener struct {
		cond    sync.Cond
		trigger uint32
		p       unsafe.Pointer
	}

	locker struct{}
)

var (
	_ sync.Locker = &locker{}
	_ Listener    = &listener{}
)

func (locker) Lock() {}

func (locker) Unlock() {}

func NewListener() Listener {
	return newListener()
}

func newListener() *listener {
	return &listener{
		cond: sync.Cond{L: locker{}},
	}
}

func (l *listener) Broadcast(value interface{}) {
	atomic.StorePointer(&l.p, unsafe.Pointer(&value))

	if atomic.CompareAndSwapUint32(&l.trigger, 0, 1) {
		l.cond.Broadcast()
	}
}

func (l *listener) Receive() (interface{}, bool) {
	if atomic.LoadUint32(&l.trigger) == 0 {
		return nil, false
	}

	p := atomic.LoadPointer(&l.p)
	if p == nil {
		return nil, true
	}

	return *(*interface{})(p), true
}

func (l *listener) Wait() interface{} {
	if atomic.LoadUint32(&l.trigger) == 0 {
		l.cond.Wait()
	}

	p := atomic.LoadPointer(&l.p)
	if p == nil {
		return nil
	}

	return *(*interface{})(p)
}
