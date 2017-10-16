package listener

import (
	"sync"
)

type (
	IntListeners struct {
		creater func() Listener
		lmap    map[int]Listener
		mu      sync.RWMutex
	}
)

func NewIntListeners(creater ...func() Listener) *IntListeners {
	var c func() Listener = NewListener
	if len(creater) != 0 && creater[0] != nil {
		c = creater[0]
	}

	return &IntListeners{
		creater: c,
		lmap:    make(map[int]Listener, 8),
	}
}

func (l *IntListeners) GetOrCreate(key int) (li Listener, found bool) {
	l.mu.RLock()
	li, found = l.lmap[key]
	l.mu.RUnlock()
	if !found {
		l.mu.Lock()
		li, found = l.lmap[key]
		if !found {
			li = l.creater()
			l.lmap[key] = li
		}
		l.mu.Unlock()
	}

	return
}

func (l *IntListeners) Get(key int) (li Listener, found bool) {
	l.mu.RLock()
	li, found = l.lmap[key]
	l.mu.RUnlock()

	return
}

func (l *IntListeners) Len() int {
	return len(l.lmap)
}

func (l *IntListeners) Delete(key int) {
	l.mu.Lock()
	delete(l.lmap, key)
	l.mu.Unlock()
}

func (l *IntListeners) Put(key int, li Listener) (old Listener) {
	l.mu.Lock()
	old = l.lmap[key]
	if li != nil {
		l.lmap[key] = li
	}
	l.mu.Unlock()

	return
}

func (l *IntListeners) Range(f func(key int, li Listener) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for key, li := range l.lmap {
		if !f(key, li) {
			break
		}
	}
}
