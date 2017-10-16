// +build go1.9

package listener

import (
	"sync"
)

type (
	Listener interface {
		Broadcast(value interface{})
		Receive() (interface{}, bool)
		Wait() interface{}
	}

	Listeners struct {
		creater func() Listener
		lmap    sync.Map //map[interface{}]Listener
	}
)

func NewListeners(creater ...func() Listener) *Listeners {
	var c func() Listener = NewListener
	if len(creater) != 0 && creater[0] != nil {
		c = creater[0]
	}

	return &Listeners{
		creater: c,
	}
}

func (l *Listeners) GetOrCreate(key interface{}) (Listener, bool) {
	li, found := l.lmap.Load(key)
	if !found {
		li = l.creater()
		li, found = l.lmap.LoadOrStore(key, li)
	}

	return li.(Listener), found
}

func (l *Listeners) Get(key interface{}) (Listener, bool) {
	if li, found := l.lmap.Load(key); found {
		return li.(Listener), true
	}

	return nil, false
}

func (l *Listeners) Len() (n int) {
	l.lmap.Range(func(key, value interface{}) bool {
		n++
		return true
	})

	return
}

func (l *Listeners) Delete(key interface{}) {
	l.lmap.Delete(key)
}

func (l *Listeners) Put(key interface{}, li Listener) Listener {
	if li == nil {
		old, _ := l.Get(key)
		return old
	}

	old, found := l.lmap.LoadOrStore(key, li)
	if !found {
		return nil
	}

	return old.(Listener)
}

func (l *Listeners) Range(f func(key interface{}, li Listener) bool) {
	l.lmap.Range(func(key interface{}, v interface{}) bool {
		return f(key, v.(Listener))
	})
}
