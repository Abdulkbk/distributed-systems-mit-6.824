package main

import (
	"fmt"
)

func main() {
	fmt.Println("=========Serial Crawler==========")
	Serial("http://golang.org/", fetcher, make(map[string]bool))

	fmt.Println("========= Concurrent Mutex =========")
	ConcurrentMutex("http://golang.org/", fetcher, makeState())

	fmt.Println("========= Concurrent Channel =========")
	ConcurrentChannel("http://golang.org/", fetcher)
}
