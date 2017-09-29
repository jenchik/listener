package listener_test

import (
	"runtime"
	"sync"
	"testing"
	"time"

	. "github.com/jenchik/listener"
	"github.com/stretchr/testify/assert"
)

func TestEmptyListeners(t *testing.T) {
	ls := NewListeners()

	li, found := ls.Get("bla bla")
	assert.Nil(t, li)
	assert.False(t, found)
	assert.Equal(t, 0, ls.Len())
}

func TestListeners(t *testing.T) {
	ls := NewListeners()

	assert.Equal(t, 0, ls.Len())

	li1, found := ls.GetOrCreate("key1")
	assert.NotNil(t, li1)
	assert.False(t, found)
	assert.Equal(t, 1, ls.Len())

	li2, found := ls.GetOrCreate("key2")
	assert.NotNil(t, li2)
	assert.False(t, found)
	assert.Equal(t, 2, ls.Len())

	liX, found := ls.GetOrCreate("key1")
	assert.True(t, liX == li1)
	assert.False(t, liX == li2)
	assert.True(t, found)
	assert.Equal(t, 2, ls.Len())

	ls.Delete("key1")
	ls.Delete("key3")
	assert.Equal(t, 1, ls.Len())

	liX, found = ls.Get("key1")
	assert.Nil(t, liX)
	assert.False(t, found)
}

func TestListenersAsync(t *testing.T) {
	start := time.Now()
	ls := NewListeners()
	var wg sync.WaitGroup

	li1, found := ls.GetOrCreate("key1")
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li1.Broadcast(123)
	})

	li2, found := ls.GetOrCreate(777)
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li2.Broadcast("foobar")
	})

	li, found := ls.Get("key3")
	assert.Nil(t, li)
	assert.False(t, found)

	value, found := li1.Receive()
	assert.False(t, found)
	assert.Nil(t, value)

	wg.Add(1)
	var v interface{}
	go func() {
		defer wg.Done()
		v = li1.Wait()
	}()

	value = li1.Wait()
	assert.Equal(t, 123, value)

	value = li2.Wait()
	assert.Equal(t, "foobar", value)

	li, found = ls.GetOrCreate(777)
	assert.NotNil(t, li)
	assert.True(t, found)
	value, found = li.Receive()
	assert.True(t, found)
	assert.Equal(t, "foobar", value)

	li1.Broadcast(312.996)
	value = li1.Wait()
	assert.Equal(t, 312.996, value)

	wg.Wait()
	assert.Equal(t, 123, v)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
	assert.Equal(t, 2, ls.Len())
}

func TestListenersBroadcast(t *testing.T) {
	start := time.Now()
	ls := NewListeners()

	type points struct {
		x int
		y int
	}

	var p1 = points{11, 22}
	var p2 = points{44, 33}
	var f = func() string {
		return "Booz!"
	}

	li1, found := ls.GetOrCreate(p1)
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li1.Broadcast(f)
	})

	li2, found := ls.GetOrCreate(&p2)
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li2.Broadcast(li1)
	})

	liX, found := ls.Get(p1)
	assert.NotNil(t, liX)
	assert.True(t, found)

	liX, found = ls.Get(p2)
	assert.Nil(t, liX)
	assert.False(t, found)

	liX, found = ls.Get(&p1)
	assert.Nil(t, liX)
	assert.False(t, found)

	liX, found = ls.Get(&p2)
	assert.NotNil(t, liX)
	assert.True(t, found)
	assert.False(t, li1 == liX)
	assert.True(t, li2 == liX)

	value := li1.Wait()
	f1, ok := value.(func() string)
	assert.True(t, ok)
	assert.Equal(t, f1(), f())

	value = li2.Wait()
	assert.Implements(t, (*Listener)(nil), value)
	liX, ok = value.(Listener)
	assert.True(t, ok)
	assert.True(t, li1 == liX)

	s := "sample"
	li1.Broadcast(&s)
	value, found = liX.Receive()
	assert.True(t, found)
	assert.Equal(t, value, &s)
	assert.Equal(t, *(value.(*string)), s)

	s = "xxyyzz"
	assert.Equal(t, *(value.(*string)), s)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestListenerOnceBroadcast(t *testing.T) {
	start := time.Now()
	ls := NewListeners(NewListenerOnce)

	s1 := "key1"
	s2 := "key1"
	ps2 := &s2

	li1, found := ls.GetOrCreate(&s1)
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li1.Broadcast("sms")
	})

	li2, found := ls.GetOrCreate(ps2)
	assert.False(t, found)
	time.AfterFunc(100*time.Millisecond, func() {
		li2.Broadcast("foo")
	})

	liX, found := ls.Get("key1")
	assert.Nil(t, liX)
	assert.False(t, found)

	liX, found = ls.Get(&s1)
	assert.NotNil(t, liX)
	assert.True(t, found)
	assert.True(t, li1 == liX)

	liX, found = ls.Get(*ps2)
	assert.Nil(t, liX)
	assert.False(t, found)

	liX, found = ls.Get(&s2)
	assert.NotNil(t, liX)
	assert.True(t, found)
	assert.True(t, li2 == liX)

	value := li1.Wait()
	assert.Equal(t, "sms", value)

	value = liX.Wait()
	assert.Equal(t, "foo", value)

	li1.Broadcast("bar")
	value = li1.Wait()
	assert.Equal(t, "sms", value)

	assert.InDelta(t, 0.1, time.Since(start).Seconds(), 0.01)
}

func TestListenersPut(t *testing.T) {
	ls := NewListeners()

	li1, found := ls.GetOrCreate("key1")
	assert.False(t, found)

	li2 := NewListener()
	assert.False(t, li2 == li1)

	li3 := NewListener()
	assert.False(t, li3 == li1)
	assert.False(t, li3 == li2)

	liX, _ := ls.Get("key2")
	assert.Nil(t, liX)

	liX, _ = ls.Get("key1")
	assert.True(t, liX == li1)
	assert.False(t, liX == li2)

	liX = ls.Put("key1", li2)
	assert.True(t, liX == li1)
	assert.False(t, liX == li2)

	liX = ls.Put("key2", li3)
	assert.Nil(t, liX)

	liX, _ = ls.Get("key2")
	assert.True(t, liX == li3)

	liX = ls.Put("key2", li3)
	assert.True(t, liX == li3)
}

func closure(l Listener) {
	s := "foobar"

	time.AfterFunc(time.Millisecond, func() {
		l.Broadcast(&s)
	})
}

type T [2][2][2][2][2][2][2][2][2][2]*int

func extFunc(l Listener) {
	a := new(T)

	a[0][0][0][0][0][0][0][0][0][0] = new(int)
	*a[0][0][0][0][0][0][0][0][0][0] = 75305

	time.AfterFunc(time.Millisecond, func() {
		l.Broadcast(a)
	})
}

func TestListenersWithGC(t *testing.T) {
	ls := NewListeners()

	li1, found := ls.GetOrCreate("key1")
	assert.False(t, found)
	li2, found := ls.GetOrCreate("key2")
	assert.False(t, found)

	closure(li1)
	extFunc(li2)

	time.Sleep(time.Millisecond * 200)
	runtime.GC()

	value := li1.Wait()
	assert.NotNil(t, value)
	pS, ok := value.(*string)
	assert.True(t, ok)
	assert.Equal(t, *pS, "foobar")

	value = li2.Wait()
	assert.NotNil(t, value)
	pT, ok := value.(*T)
	assert.True(t, ok)
	assert.Equal(t, *pT[0][0][0][0][0][0][0][0][0][0], 75305)
}
