package main

import (
	"fmt"
	"sync"
)

// Fetcher
type Fetcher interface {
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeResult is a struct to store the body and urls of a page.
type fakeResult struct {
	body string
	urls []string
}

// fakeFetcher is a Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found: %s\n", url)
		return res.urls, nil
	}

	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

// Serial crawls a website serially.
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	// Return early if we already have fetched the url.
	if fetched[url] {
		return
	}

	fetched[url] = true

	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}

	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
}

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
