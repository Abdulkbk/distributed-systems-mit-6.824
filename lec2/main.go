package main

import (
	"fmt"
	"sync"
)

// Concurrent Crawler with shared state and Mutex.
type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

// makeState returns fetchState
func makeState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)

	return f
}

// ConcurrentMutex crawls a website concurrently using mutex to protect shared state.
func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}

	var wg sync.WaitGroup

	for _, u := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u)
	}

	wg.Wait()
}

//
// Concurrent crawler with Channels
//

// worker
func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		ch <- []string{}
	} else {
		ch <- urls
	}
}

// master
func master(ch chan []string, fetcher Fetcher) {
	n := 1

	fetched := make(map[string]bool)

	for urls := range ch {
		for _, u := range urls {
			if fetched[u] == false {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher)
			}
		}

		n -= 1

		if n == 0 {
			break
		}
	}

}

// ConcurrentChannel
func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url}
	}()

	master(ch, fetcher)

}

func main() {
	fmt.Println("=========Serial Crawler==========")
	Serial("http://golang.org/", fetcher, make(map[string]bool))

	fmt.Println("========= Concurrent Mutex =========")
	ConcurrentMutex("http://golang.org/", fetcher, makeState())

	fmt.Println("========= Concurrent Channel =========")
	ConcurrentChannel("http://golang.org/", fetcher)
}
