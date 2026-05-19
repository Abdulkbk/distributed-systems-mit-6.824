package main

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
