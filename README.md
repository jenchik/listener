# listener
Concurent Sub/Unsub pattern for Golang. Can be used as an alternative to channels

[![Build Status](https://travis-ci.org/jenchik/listener.svg)](https://travis-ci.org/jenchik/listener)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenchik/listener)](https://goreportcard.com/report/github.com/jenchik/listener)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjenchik%2Flistener.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjenchik%2Flistener?ref=badge_shield)

Installation
------------

```bash
go get github.com/jenchik/listener
```

Examples
-------

with waiting
```go
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jenchik/listener"
)

func sendRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func worker(l listener.Listener, url string) {
	response, err := sendRequest(url)
	if err != nil {
		fmt.Println("Error request:", err)
		return
	}

	if response.Body == nil {
		fmt.Println("Error empty body")
		return
	}
	defer response.Body.Close()

	l.Broadcast(response.Status)

	// ...
}

func main() {
	start := time.Now()

	pool := listener.NewListeners()

	var wg sync.WaitGroup
	wg.Add(5)
	urlMask := "http://example.com/page/%d"
	for i := 0; i < 5; i++ {
		l, _ := pool.GetOrCreate(i)
		go worker(l, fmt.Sprintf(urlMask, i))
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Request with ID-->%d is '%s'\n", id, l.Wait())
		}(i)
	}

	// ... other work

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second * 2):
		fmt.Println("Error with timeout")
	}

	fmt.Println("Duration", time.Since(start).String())
}
```

without waiting
```go
package main

import (
	"fmt"
	"time"

	"github.com/jenchik/listener"
)

type (
	result struct {
		id int
		n  uint64
	}
)

func (r *result) String() string {
	return fmt.Sprintf("Worker ID-->%d; Result-->%d", r.id, r.n)
}

func worker(response, stop listener.Listener, id int) {
	sum := uint64(id * 1000)
	tic := time.NewTicker(time.Millisecond)

	for {
		select {
		case <-tic.C:
			response.Broadcast(&result{id, sum})
		default:
			if _, ok := stop.Receive(); ok {
				tic.Stop()
				return
			}
		}
		sum++
	}
}

func dispatcher(pool *listener.Listeners, count int) {
	stop, _ := pool.GetOrCreate("stop")
	for i := 0; i < count; i++ {
		l, _ := pool.GetOrCreate(i)
		go worker(l, stop, i)
	}

	for {
		if _, ok := stop.Receive(); ok {
			return
		}
		<-time.After(time.Second)
		if _, ok := stop.Receive(); ok {
			return
		}
		for i := 0; i < count; i++ {
			l, _ := pool.GetOrCreate(i)
			if v, ok := l.Receive(); ok {
				fmt.Println(v)
			}
		}
		fmt.Println("")
	}
}

func main() {
	start := time.Now()

	pool := listener.NewListeners()

	go dispatcher(pool, 10)

	time.Sleep(time.Second * 10)

	stop, _ := pool.GetOrCreate("stop")
	stop.Broadcast(false)

	fmt.Println("Duration", time.Since(start).String())
}
```

Benchmarks
----------
```
BenchmarkThreadsResend-4             200000000             9.65 ns/op     207.20 MB/s           0 B/op           0 allocs/op
BenchmarkThreadsResendString-4       200000000             9.30 ns/op     215.14 MB/s           0 B/op           0 allocs/op
BenchmarkThreadsResendInt-4          200000000             9.14 ns/op     218.80 MB/s           0 B/op           0 allocs/op
BenchmarkThreadsOnce-4               200000000             9.96 ns/op     200.79 MB/s           0 B/op           0 allocs/op
BenchmarkThreadsOnceString-4         200000000             9.78 ns/op     204.54 MB/s           0 B/op           0 allocs/op
BenchmarkThreadsOnceInt-4            200000000             9.21 ns/op     217.20 MB/s           0 B/op           0 allocs/op
```


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fjenchik%2Flistener.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fjenchik%2Flistener?ref=badge_large)