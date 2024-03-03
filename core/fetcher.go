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
		summary:   make(map[StatusCodeGroup]int),
	}
}

func (f *Fetcher) Summary() {
	success := f.summary[SUCCESSFUL_RES]
	var fails int
	for key, val := range f.summary {
		if key != SUCCESSFUL_RES {
			fails += val
			continue
		}
	}
	fmt.Println("Successes: ", success)
	fmt.Println("Failures: ", fails)
}

func (f *Fetcher) Display() {
	for _, response := range f.responses {
		fmt.Println("Response code: ", response.Code)
	}
}

func (f *Fetcher) Run() {
	switch f.fetchType {
	case SEQUENTIAL:
		f.sequenceFetching()
	case CONCURRENT:
		responses := f.concurrentFetching(f.url, f.quantity, f.limit)
		f.responses = responses
	}
}

func (f *Fetcher) fetch() Response {
	resp, err := http.Get(f.url)
	if err != nil {
		log.Println(err)
		return Response{}
	}
	defer resp.Body.Close()
	return Response{Code: resp.StatusCode}
}

func (f *Fetcher) sequenceFetching() {
	f.responses = make([]Response, 0, f.quantity)
	for i := 0; i < f.quantity; i++ {
		response := f.fetch()
		f.countStatus(response.Code)
		f.responses = append(f.responses, response)
	}
}

func (f *Fetcher) concurrentFetching(url string, number int, maxConcurrent int) []Response {
	results := make([]Response, 0, number)
	message := make(chan Response, number)
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	for i := range number {
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() {
				wg.Done()
				<-sem
			}()
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
				message <- Response{}
				return
			}
			defer resp.Body.Close()

			message <- Response{Code: resp.StatusCode}
		}()

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
		results = append(results, result)
	}

	return results
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
	Code int
}
