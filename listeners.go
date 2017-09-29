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
		lmap    map[interface{}]Listener
		mu      sync.Mutex
	}
)

func NewListeners(creater ...func() Listener) *Listeners {
	var c func() Listener = NewListener
	if len(creater) != 0 && creater[0] != nil {
		c = creater[0]
	}

	return &Listeners{
		creater: c,
		lmap:    make(map[interface{}]Listener, 8),
	}
}

func (l *Listeners) GetOrCreate(key interface{}) (li Listener, found bool) {
	l.mu.Lock()
	li, found = l.lmap[key]
	if !found {
		li = l.creater()
		l.lmap[key] = li
	}
	l.mu.Unlock()

	return
}

func (l *Listeners) Get(key interface{}) (li Listener, found bool) {
	l.mu.Lock()
	li, found = l.lmap[key]
	l.mu.Unlock()

	return
}

func (l *Listeners) Len() int {
	return len(l.lmap)
}

func (l *Listeners) Delete(key interface{}) {
	l.mu.Lock()
	delete(l.lmap, key)
	l.mu.Unlock()
}

func (l *Listeners) Put(key interface{}, li Listener) (old Listener) {
	l.mu.Lock()
	old = l.lmap[key]
	if li != nil {
		l.lmap[key] = li
	}
	l.mu.Unlock()

	return
}
