package core

import (
	"strings"
	"testing"
)

// Assume the server is already started at port 8000
// To start the mock server, please go to directory www/
// and run python -m http.server

const TEST_URL = "http://0.0.0.0:8000"

func TestFetcherRun(t *testing.T) {
	request := 100
	concurrent := 1

	fetcher := NewFetcher(TEST_URL, request, concurrent)
	fetcher.Run()
	output := fetcher.genSummary(fetcher.calcSummary())

	// Ensure Summary Header output is valid
	if !strings.Contains(output, "Results") {
		t.Errorf("Unexpected Summary Header: \n%s", output)
	}

	// Ensure Total Request output is valid
	if !strings.Contains(output, "Total Requests") {
		t.Errorf("Expected 'Total Request' string but not exists: \n%s", output)
	}

	// Ensure Failed Request output is valid
	if !strings.Contains(output, "Failed Requests") {
		t.Errorf("Expected 'Failed Requests' string but not exists: \n%s", output)
	}
}

func BenchmarkFetcherRun(b *testing.B) {
	request := 10000
	concurrent := 1

	fetcher := NewFetcher(TEST_URL, request, concurrent)
	fetcher.Run()
}

func BenchmarkFetcherRunConcurrent(b *testing.B) {
	request := 10000
	concurrent := 10

	fetcher := NewFetcher(TEST_URL, request, concurrent)
	fetcher.Run()
}
