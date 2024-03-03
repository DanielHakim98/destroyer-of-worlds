package core

import (
	"fmt"
	"log"
	"net/http"
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

type Fetcher struct {
	url       string
	quantity  int
	responses []Response
	summary   map[StatusCodeGroup]int
}

func NewFetcher(url string, n int) *Fetcher {
	return &Fetcher{
		url:      url,
		quantity: n,
		summary:  make(map[StatusCodeGroup]int),
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
	f.responses = make([]Response, 0, f.quantity)
	for i := 0; i < f.quantity; i++ {
		response := f.fetch()
		f.countStatus(response.Code)
		f.responses = append(f.responses, response)
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
