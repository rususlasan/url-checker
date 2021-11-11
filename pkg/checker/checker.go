package checker

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type urlChecker struct {
	client *http.Client

	concurrency int // how many workers to run
}

type result struct {
	url      string
	bodySize int64
	err      error
}

// NewChecker creates new instance of urlChecker, timeout is a timeout for http requests
func NewChecker(concurrency int, timeout time.Duration) (*urlChecker, error) {
	if int64(concurrency)+timeout.Nanoseconds() < 2 {
		return nil, errors.New("wrong parameters for checker")
	}

	c := &urlChecker{
		concurrency: concurrency,
	}

	c.client = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	return c, nil
}

func (c *urlChecker) Check(urls []string) {
	var (
		globalWg    = &sync.WaitGroup{}
		workersWg   = &sync.WaitGroup{}
		inChan      = make(chan string, c.concurrency)
		resultsChan = make(chan *result, c.concurrency)
	)

	globalWg.Add(1)
	// feed inChan
	go func() {
		globalWg.Done()
		for _, url := range urls {
			inChan <- url
		}
		close(inChan)
	}()

	// run workers
	workersWg.Add(c.concurrency)
	for i := 0; i < c.concurrency; i++ {
		go c.worker(inChan, resultsChan, workersWg)
	}

	// build summary response
	globalWg.Add(1)
	go func() {
		defer globalWg.Done()
		successResults := make([]*result, 0, len(urls))
		for res := range resultsChan {
			if res.err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: check for url '%s' failed with: %s\n", res.url, res.err)
				continue
			}
			successResults = append(successResults, res)
		}

		sort.Slice(successResults, func(i, j int) bool {
			return successResults[i].bodySize < successResults[j].bodySize
		})

		for _, res := range successResults {
			fmt.Fprintf(os.Stdout, "%s %d\n", res.url, res.bodySize)
		}
	}()

	workersWg.Wait()
	close(resultsChan)
	globalWg.Wait()
}

// worker reads from in chan, checks url and calculates body size then sends results to the resChan
func (c *urlChecker) worker(inChan <-chan string, resultsChan chan<- *result, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		resp *http.Response
		err  error
	)
	for url := range inChan {
		res := result{url: url}
		resp, err = c.client.Get(url)
		if err != nil {
			res.err = err
		} else if resp != nil && resp.Body != nil {
			res.bodySize, res.err = io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}

		resultsChan <- &res
	}
}
