package listener

import (
	"sync"
)

type (
	StringListeners struct {
		creater func() Listener
		lmap    map[string]Listener
		mu      sync.RWMutex
	}
)

func NewStringListeners(creater ...func() Listener) *StringListeners {
	var c func() Listener = NewListener
	if len(creater) != 0 && creater[0] != nil {
		c = creater[0]
	}

	return &StringListeners{
		creater: c,
		lmap:    make(map[string]Listener, 8),
	}
}

func (l *StringListeners) GetOrCreate(key string) (li Listener, found bool) {
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

func (l *StringListeners) Get(key string) (li Listener, found bool) {
	l.mu.RLock()
	li, found = l.lmap[key]
	l.mu.RUnlock()

	return
}

func (l *StringListeners) Len() int {
	return len(l.lmap)
}

func (l *StringListeners) Delete(key string) {
	l.mu.Lock()
	delete(l.lmap, key)
	l.mu.Unlock()
}

func (l *StringListeners) Put(key string, li Listener) (old Listener) {
	l.mu.Lock()
	old = l.lmap[key]
	if li != nil {
		l.lmap[key] = li
	}
	l.mu.Unlock()

	return
}

func (l *StringListeners) Range(f func(key string, li Listener) bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for key, li := range l.lmap {
		if !f(key, li) {
			break
		}
	}
}
