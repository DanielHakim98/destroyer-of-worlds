package core

import (
	"log"
	"net/http"
)

type Fetcher struct {
	url string
}

func NewFetcher(url string) *Fetcher {
	return &Fetcher{
		url: url,
	}
}

func (f *Fetcher) Fetch() Response {
	resp, err := http.Get(f.url)
	if err != nil {
		log.Println(err)
		return Response{}
	}
	defer resp.Body.Close()
	return Response{Code: resp.StatusCode}
}

type Response struct {
	Code int
}
