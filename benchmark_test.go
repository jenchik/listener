package listener

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

const (
	chars  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	items  = 10000
	lenKey = 10
	steps  = 500

	benchWorks = 1000
	testWorks  = 1000
)

var (
	keys        []string
	newKeys     []string
	dispersion  []string
	dispersions [][]string
)

func init() {
	keys = make([]string, 0, items)
	newKeys = make([]string, 0, items)
	var key string
	for i := 0; i < items; i++ {
		key = randString(lenKey)
		keys = append(keys, key)
		newKeys = append(newKeys, randString(lenKey))
	}

	dispersion = newDispersion(dispersion, steps)
	for i := 0; i < benchWorks; i++ {
		s := make([]string, 0, steps)
		dispersions = append(dispersions, newDispersion(s, steps))
	}
}

func randString(n int) string {
	buf := make([]byte, n)
	l := len(chars)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < n; i++ {
		buf[i] = chars[rand.Intn(l)]
	}
	return string(buf)
}

func newDispersion(in []string, x int) []string {
	rand.Seed(time.Now().UTC().UnixNano())
	var key string
	var n int
	for i := 0; i < x; i++ {
		n = rand.Intn(items)
		if rand.Intn(2) == 1 {
			key = newKeys[n]
		} else {
			key = keys[n]
		}
		in = append(in, key)
	}
	return in
}

func initMap() map[string]int {
	m := make(map[string]int, items*2)
	for i, key := range keys {
		m[key] = i
	}
	return m
}

func BenchmarkResend(b *testing.B) {
	m := initMap()
	obs := NewListeners()
	var found bool
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k string) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(key)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(key)
		}
	}
}

func BenchmarkResendString(b *testing.B) {
	m := initMap()
	obs := NewStringListeners()
	var found bool
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k string) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(key)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(key)
		}
	}
}

func BenchmarkResendInt(b *testing.B) {
	m := initMap()
	obs := NewIntListeners()
	var found bool
	var key string
	var keyInt int

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for keyInt, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k int) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(k)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(keyInt)
		}
	}
}

func BenchmarkOnce(b *testing.B) {
	m := initMap()
	obs := NewListeners(NewListenerOnce)
	var found bool
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k string) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(key)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(key)
		}
	}
}

func BenchmarkOnceString(b *testing.B) {
	m := initMap()
	obs := NewStringListeners(NewListenerOnce)
	var found bool
	var key string

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k string) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(key)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(key)
		}
	}
}

func BenchmarkOnceInt(b *testing.B) {
	m := initMap()
	obs := NewIntListeners(NewListenerOnce)
	var found bool
	var key string
	var keyInt int

	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for keyInt, key = range dispersion {
			if _, found = m[key]; found {
				continue
			}

			go func(k int) {
				l, f := obs.GetOrCreate(k)
				if !f {
					time.AfterFunc(time.Millisecond, func() {
						obs.Delete(k)
						l.Broadcast(312)
					})
				}
				if l.Wait().(int) != 312 {
					b.Fail()
				}
			}(keyInt)
		}
	}
}

func BenchmarkThreadsResend(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewListeners()

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var l Listener
		var i int
		for pb.Next() {
			key = disp[i%steps]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(key)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(key)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}

func BenchmarkThreadsResendString(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewStringListeners()

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var l Listener
		var i int
		for pb.Next() {
			key = disp[i%steps]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(key)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(key)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}

func BenchmarkThreadsResendInt(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewIntListeners()

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var keyInt int
		var l Listener
		var i int
		for pb.Next() {
			keyInt = i % steps
			key = disp[keyInt]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(keyInt)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(keyInt)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}

func BenchmarkThreadsOnce(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewListeners(NewListenerOnce)

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var l Listener
		var i int
		for pb.Next() {
			key = disp[i%steps]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(key)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(key)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}

func BenchmarkThreadsOnceString(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewStringListeners(NewListenerOnce)

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var l Listener
		var i int
		for pb.Next() {
			key = disp[i%steps]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(key)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(key)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}

func BenchmarkThreadsOnceInt(b *testing.B) {
	var d uint32

	m := initMap()
	obs := NewIntListeners(NewListenerOnce)

	b.SetParallelism(benchWorks)
	b.ReportAllocs()
	b.SetBytes(2)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		dd := atomic.AddUint32(&d, 1)
		disp := dispersions[int(dd)%benchWorks]
		var found bool
		var key string
		var keyInt int
		var l Listener
		var i int
		for pb.Next() {
			keyInt = i % steps
			key = disp[keyInt]
			if _, found = m[key]; found {
				continue
			}

			l, found = obs.GetOrCreate(keyInt)
			if !found {
				time.AfterFunc(time.Millisecond, func() {
					obs.Delete(keyInt)
					l.Broadcast(312)
				})
			}
			if l.Wait().(int) != 312 {
				b.Fail()
			}

			i++
		}
	})
}
