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

const (
	SUMMARY_HEADER_DISPLAY       = "Results: "
	TOTAL_REQUEST_DISPLAY        = "\n Total Requests  (2XX)                        .......................: "
	FAILED_REQUEST_DISPLAY       = "\n Failed Requests (4XX, 5XX and unknown)       .......................: "
	TOTAL_EXECUTION_TIME_DISPLAY = "\n Total execution time                         .......................: "
	TOTAL_REQUESTS_TIME_DISPLAY  = "\n Total requests  time                         .......................: "
	REQUEST_PER_SECOND_DISPLAY   = "\n Request/second                               .......................: "
	STATISTIC_HEADER_DISPLAY     = "\n\nStatistic:"
	REQUEST_TIME_STATS_DISPLAY   = "\n Request Time (s) (Min, Max, Mean)            .......................: "
	STATS_TEMPL                  = "(%v, %v, %v)\n"
)

type Fetcher struct {
	url       string
	quantity  int
	limit     int
	fetchType FetchType
	responses []Response
	execTime  time.Duration
	summary   map[StatusCodeGroup]int
	stats     Stats
}

type Stats struct {
	statusCodes  map[StatusCodeGroup]int
	requestsTime [3]time.Duration
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
		stats: Stats{
			statusCodes: make(map[StatusCodeGroup]int),
		},
	}
}

func (f *Fetcher) Summary() {
	summaryStr := f.genSummary(f.calcSummary())
	fmt.Println(summaryStr)
}

func (f *Fetcher) calcSummary() Summary {
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

	min := f.stats.requestsTime[0]
	max := f.stats.requestsTime[1]
	mean := totalDuration.Seconds() / total

	return Summary{
		total:         total,
		average:       average,
		execTime:      execTime,
		totalDuration: totalDuration,
		fails:         fails,
		totalDurStats: struct {
			min  time.Duration
			max  time.Duration
			mean float64
		}{
			min, max, mean,
		},
	}
}

type Summary struct {
	total, average          float64
	execTime, totalDuration time.Duration
	fails                   int
	totalDurStats           struct {
		min, max time.Duration
		mean     float64
	}
}

func (f *Fetcher) genSummary(s Summary) string {
	return SUMMARY_HEADER_DISPLAY +
		TOTAL_REQUEST_DISPLAY + fmt.Sprint(s.total) +
		FAILED_REQUEST_DISPLAY + fmt.Sprint(s.fails) +
		TOTAL_EXECUTION_TIME_DISPLAY + s.execTime.String() +
		TOTAL_REQUESTS_TIME_DISPLAY + s.totalDuration.String() +
		REQUEST_PER_SECOND_DISPLAY + fmt.Sprint(s.average) +
		STATISTIC_HEADER_DISPLAY +
		REQUEST_TIME_STATS_DISPLAY + fmt.Sprintf(
		STATS_TEMPL, s.totalDurStats.min.Seconds(), s.totalDurStats.max.Seconds(), s.totalDurStats.mean)
}

func (f *Fetcher) Display() {
	for _, response := range f.responses {
		fmt.Println("Response code: ", response.Code)
	}
}

func (f *Fetcher) Run() {
	switch f.fetchType {
	case SEQUENTIAL:
		timer := logTime()
		f.sequenceFetching()
		f.execTime = timer()
	case CONCURRENT:
		timer := logTime()
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
		f.genStats(response)
		f.responses = append(f.responses, response)
	}
}

func (f *Fetcher) genStats(res Response) {
	f.countStatusCode(res.Code)
	f.findMaxMinDur(res.Duration)
}

func (f *Fetcher) findMaxMinDur(t time.Duration) {
	curMin := f.stats.requestsTime[0]
	// Only if t is smaller then current min, then replace
	if curMin == 0 || t < curMin {
		f.stats.requestsTime[0] = t
	}

	curMax := f.stats.requestsTime[1]
	// If zero valued or larger than current max, then replace
	if curMax == 0 || t > curMax {
		f.stats.requestsTime[1] = t
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
		f.genStats(result)
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

func (f *Fetcher) countStatusCode(code int) {
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

func logTime() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		execTime := time.Since(start)
		return execTime
	}
}
