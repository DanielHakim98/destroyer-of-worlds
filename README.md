# Destroyer of Worlds

- It's just a normal load tester, I build it because I'm trying this challenge **[coding-challenges: load-tester](https://codingchallenges.fyi/challenges/challenge-load-tester/)**.

## Installation

- For now, you need to build it from source. **Golang 1.22** is required to build the binary. **Python 3.12.0** is used to set up mock server serves static html.

## Usage

```text
Usage:
destroyer-of-worlds [flags]

Flags:
-c, --concurrent int maximum concurent request. Default is 1 (default 1)
-h, --help help for destroyer-of-worlds
-n, --requests int The total requests to be sent. Default is 1 (default 1)
-t, --toggle Help message for toggle
-u, --url string URL to be tested.

```

To setup mock server, make sure current shell is at root project directory, and then run command below

```text
cd www
python -m http.server
```

## Issue

- Refer details in core/fetcher_test.go, I believe the concurrency code is kind of suboptimal. This is the result of benchmarking runs locally and the server tested is python mock server, also run locally:

  ```text
  $ go test ./core/ -bench=.

  goos: linux
  goarch: amd64
  pkg: github.com/DanielHakim98/destroyer-of-worlds/core
  cpu: Intel(R) Core(TM) i5-8300H CPU @ 2.30GHz
  BenchmarkFetcherRun-8                          1        5947048191 ns/op
  BenchmarkFetcherRunConcurrent-8                1        51590994485 ns/op
  PASS
  ok      github.com/DanielHakim98/destroyer-of-worlds/core       57.676s
  ```

  I suspect there might be overhead when setting up concurrency and causes. Will update later after profiling (If I'm not lazy lah)

## Contributing

Nah, I don't think it's worth it anyway. You can just copy and paste, steal it, or even recreate to do something better, robust and more interesting that this. Or just use **[Hey](https://github.com/rakyll/hey)**.
