package core

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type StatusCodeGroup int

const (
	INFORMATION_RES StatusCodeGroup = iota + 1
	SUCCESSFUL_RES
	REDIRECTION_RES
	CLIENT_ERROR_RES
	SERVER_ERROR_RES
	UNKNOWN_RES
)

type FetchType int

const (
	SEQUENTIAL FetchType = iota
	CONCURRENT
)

type Fetcher struct {
	url       string
	quantity  int
	limit     int
	fetchType FetchType
	responses []Response
	execTime  time.Duration
	summary   map[StatusCodeGroup]int
}

func NewFetcher(url string, number int, concurrent int) *Fetcher {
	var t FetchType
	if concurrent == 1 {
		t = SEQUENTIAL
	} else {
		t = CONCURRENT
	}

	return &Fetcher{
		url:       url,
		quantity:  number,
		limit:     concurrent,
		fetchType: t,
		responses: make([]Response, 0, number),
		summary:   make(map[StatusCodeGroup]int),
	}
}

func (f *Fetcher) Summary() {
	total := float64(f.quantity)
	var fails int
	for key, val := range f.summary {
		if key == CLIENT_ERROR_RES || key == SERVER_ERROR_RES || key == UNKNOWN_RES {
			fails += val
			continue
		}
	}

	var totalDuration time.Duration
	for _, response := range f.responses {
		totalDuration += response.Duration
	}
	average := total / totalDuration.Seconds()

	execTime := f.execTime

	fmt.Println("Results: ")
	fmt.Println(" Total Requests  (2XX)                  .......................: ", total)
	fmt.Println(" Failed Requests (4XX, 5XX and unknown) .......................: ", fails)
	fmt.Println(" Total execution time                   .......................: ", execTime)
	fmt.Println(" Total requests  time                   .......................: ", totalDuration)
	fmt.Println(" Request/second                         .......................: ", average)
}

func (f *Fetcher) Display() {
	for _, response := range f.responses {
		fmt.Println("Response code: ", response.Code)
	}
}

func (f *Fetcher) Run() {
	switch f.fetchType {
	case SEQUENTIAL:
		timer := logTime("sequenceFetching")
		f.sequenceFetching()
		f.execTime = timer()
	case CONCURRENT:
		timer := logTime("concurrentFetching")
		f.concurrentFetching(f.url, f.quantity, f.limit)
		f.execTime = timer()
	}
}

func (f *Fetcher) fetch() Response {
	start := time.Now()
	resp, err := http.Get(f.url)
	if err != nil {
		log.Println(err)
		return Response{
			Duration: time.Since(start),
		}
	}
	defer resp.Body.Close()
	return Response{
		Code:     resp.StatusCode,
		Duration: time.Since(start),
	}
}

func (f *Fetcher) sequenceFetching() {
	for i := 0; i < f.quantity; i++ {
		response := f.fetch()
		f.countStatus(response.Code)
		f.responses = append(f.responses, response)
	}
}

func (f *Fetcher) concurrentFetching(url string, number int, maxConcurrent int) {
	message := make(chan Response, number)
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	for i := range number {
		wg.Add(1)
		go f.worker(&wg, sem, url, message)
		if i%maxConcurrent == 0 && i != number {
			time.Sleep(50 * time.Millisecond)
		}
	}

	go func() {
		wg.Wait()
		close(message)
	}()

	for result := range message {
		f.countStatus(result.Code)
		f.responses = append(f.responses, result)
	}
}

func (f *Fetcher) worker(wg *sync.WaitGroup, sem chan struct{}, url string, message chan<- Response) {
	sem <- struct{}{}
	defer func() {
		wg.Done()
		<-sem
	}()

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		message <- Response{
			Duration: time.Since(start),
		}
		return
	}
	defer resp.Body.Close()

	message <- Response{
		Code:     resp.StatusCode,
		Duration: time.Since(start),
	}
}

func (f *Fetcher) countStatus(code int) {
	helper := func(group StatusCodeGroup) {
		_, exists := f.summary[group]
		if exists {
			f.summary[group]++
			return
		}
		f.summary[group] = 1
	}

	switch {
	case code >= 200 && code <= 299:
		helper(SUCCESSFUL_RES)
	case code >= 100 && code <= 199:
		helper(INFORMATION_RES)
	case code >= 300 && code <= 399:
		helper(REDIRECTION_RES)
	case code >= 400 && code <= 499:
		helper(CLIENT_ERROR_RES)
	case code >= 500 && code <= 599:
		helper(SERVER_ERROR_RES)
	default:
		helper(UNKNOWN_RES)
	}
}

type Response struct {
	Code     int
	Duration time.Duration
}

func logTime(name string) func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		execTime := time.Since(start)
		fmt.Printf("'%v' executtion time is %v\n", name, execTime)
		return execTime
	}
}
