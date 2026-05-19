package main

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
