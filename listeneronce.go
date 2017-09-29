package listener

type (
	listenerOnce struct {
		done  chan struct{}
		value interface{}
	}
)

var (
	_ Listener = &listenerOnce{}
)

func NewListenerOnce() Listener {
	return newListenerOnce()
}

func newListenerOnce() *listenerOnce {
	return &listenerOnce{
		done: make(chan struct{}),
	}
}

func (l *listenerOnce) Broadcast(value interface{}) {
	select {
	case <-l.done:
	default:
		l.value = value
		close(l.done)
	}
}

func (l *listenerOnce) Receive() (interface{}, bool) {
	select {
	case <-l.done:
	default:
		return nil, false
	}

	return l.value, true
}

func (l *listenerOnce) Wait() interface{} {
	<-l.done

	return l.value
}
