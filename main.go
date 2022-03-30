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

	tResultValue := 0

	tClient := &http.Client{
		Timeout: time.Duration(inRequestTimeout) * time.Second,
	}

	tResponse, tErr := tClient.Get(inURL)
	if tErr != nil {
		return tResultValue
	}
	defer tResponse.Body.Close()

	tContent, tErr := ioutil.ReadAll(tResponse.Body)
	if tErr != nil {
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

	tActiveJobs := make(chan struct{}, inPoolSize)
	tResultChan := make(chan int, len(inURLList)) //
	var wg sync.WaitGroup
	tTotalCount := 0

	// workers runner
	for tJobIndex := 0; tJobIndex < len(inURLList); tJobIndex++ {

		// add single job (if possible - main limiter will wait for channel slot free)
		tActiveJobs <- struct{}{}
		wg.Add(1)

		go func(inWord string, inURL string, inRequestTimeout int) {
			defer wg.Done()

			tResult := fCountWordsWorker(inWord, inURL, inRequestTimeout) // silent version - any errors just ignored (unsafe?) and returning default value
			tResultChan <- tResult
			fmt.Printf("Count for %s: %d\n", inURL, tResult)

			// free job slot
			<-tActiveJobs
		}(inWord, inURLList[tJobIndex], inRequestTimeout)
	}

	// finalyze
	wg.Wait()

	close(tResultChan)
	for tValue := range tResultChan {
		tTotalCount += tValue
	}

	close(tActiveJobs)
	fmt.Printf("Total: %d\n", tTotalCount)
}

// main
func main() {

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
