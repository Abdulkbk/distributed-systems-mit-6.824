package main

import "sync"

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
