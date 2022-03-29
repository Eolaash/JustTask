package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// fCountWordsWorker - parse web page by inURL for inWord in body content ()
// REQ: caching will speed up such request (but not all content static in time) or FILTERING base URLList to remove duplicates
func fCountWordsWorker(inWord string, inURL string, inRequestTimeout int) int {
	// set default value
	var tResultValue = 0

	// create HTTP client with timeout
	tClient := &http.Client{
		Timeout: time.Duration(inRequestTimeout) * time.Second,
	}

	// make GET request with preseted timeout
	tResponse, tErr := tClient.Get(inURL)
	if tErr != nil {
		// log.Fatal(tErr)
		return tResultValue
	}

	// responce should be closed anyway later
	defer tResponse.Body.Close()

	// read responce body as []byte (can be converted to string later)
	tContent, tErr := ioutil.ReadAll(tResponse.Body)
	if tErr != nil {
		// log.Fatal(tErr)
		return tResultValue
	}

	// get count instrings in content (how can CODING affect result? need tests)
	tResultValue = strings.Count(string(tContent), inWord)

	return tResultValue
}

// fCountWords - get inWord count in web page content listed in inURLList (unsing goroutines with parallel limitaion - inPoolSize)
// inRequestTimeout - timeout in seconds for GET requests
// REQ: May be better to make universal (goroutine)-wrapper for limited goroutines run as object or something like (as simplified pythons thread pool executor)
func fCountWords(inWord string, inURLList []string, inPoolSize int, inRequestTimeout int) {

	//checkers ...
	if inPoolSize < 1 || inRequestTimeout < 1 {
		return
	}

	// channel for limiting parallel goroutines (channel in role of slots)
	tActiveJobs := make(chan struct{}, inPoolSize)

	// syncronyzer (for all jobs to be done)
	var wg sync.WaitGroup

	// mutex for avoiding any collisions on shared resource (tTotalCount) ...oversafed?
	var tLocker sync.Mutex

	// total value default (int)
	tTotalCount := 0

	// workers runner
	for tJobIndex := 0; tJobIndex < len(inURLList); tJobIndex++ {

		// add single job (if possible - main limiter will wait for channel slot free)
		tActiveJobs <- struct{}{}
		wg.Add(1)

		// anonymous goroutine call with job remover after it (wraping worker)
		go func(inWord string, inURL string, inRequestTimeout int) {
			defer wg.Done()

			tResult := fCountWordsWorker(inWord, inURL, inRequestTimeout) // silent version - any errors just ignored (unsafe?) and returning default value

			// using mutex to work with shared variable (error with operation can create deadlock?)
			tLocker.Lock()
			tTotalCount += tResult // safe modification under mutex locker (is unbuffered channel will be better than mutex?)
			tLocker.Unlock()

			fmt.Printf("Count for %s: %d\n", inURL, tResult)

			// free job slot
			<-tActiveJobs
		}(inWord, inURLList[tJobIndex], inRequestTimeout) // passing parameters
	}

	// waiting any active gouroutines
	wg.Wait()

	// Close the channel
	close(tActiveJobs)

	// Show results
	fmt.Printf("Total: %d\n", tTotalCount)
}

// main
func main() {

	// presets
	tMaxWorkers := 5
	tTargetWord := "Go"   // case sensetive
	tRequestTimeout := 30 // request timeout in seconds
	tURLList := []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.example.com/",
		"https://dev.to/",
		"http://www.typescriptlang.org/",
		"http://www.japan.com/",
		"http://metanit.com/",
		"https://go.dev/",
		"http://www.golang.org/",
	}

	fCountWords(tTargetWord, tURLList, tMaxWorkers, tRequestTimeout)
}
