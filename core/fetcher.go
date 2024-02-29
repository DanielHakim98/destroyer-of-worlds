package core

import (
	"fmt"
	"log"
	"net/http"
)

type Fetcher struct {
	url      string
	quantity int
}

func NewFetcher(url string, n int) *Fetcher {
	return &Fetcher{
		url:      url,
		quantity: n,
	}
}

func (f *Fetcher) Run() {
	responses := make([]Response, 0, f.quantity)
	for i := 0; i < f.quantity; i++ {
		response := f.fetch()
		responses = append(responses, response)
	}

	f.display(&responses)
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

func (f *Fetcher) display(responses *[]Response) {
	for _, response := range *responses {
		fmt.Println("Response code: ", response.Code)
	}
}

type Response struct {
	Code int
}
