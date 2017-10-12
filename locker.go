package listener

import (
	"sync"
)

type (
	locker struct{}
)

var (
	_ sync.Locker = &locker{}
)

func (locker) Lock() {}

func (locker) Unlock() {}
