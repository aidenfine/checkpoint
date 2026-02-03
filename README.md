# Checkpoint


<div align="center">


![Checkpoint logo](./assets/logo.svg)
</div>

## Installation

`go get -u github.com/aidenfine/checkpoint`

## Quick Start
You can use checkpoint in two methods.

No config setup

```go
import (
    checkpoint "github.com/aidenfine/checkpoint"
)
func main() {
    rateLimiterMiddleware := checkpoint.LimitByIp(25, 1, 1)
    httpHandler := rateLimiterMiddleware(http.DefaultServeMux)

    http.ListenAndServe(":8080", httpHandler)

}
```

Or with a config for people who require a little more customization.

```go
import (
    checkpoint "github.com/aidenfine/checkpoint"
)
func main() {

	config := checkpoint.Config{
		IgnorePaths:     []string{"/logs"},
		MaxTokens:       25,
		RefillRate:      1,
		TokensPerRefill: 1,
		LimitMethod:     checkpoint.LimitByIp,
	}

    rateLimiterMiddleware := checkpoint.WithConfig(config)
    httpHandler := rateLimiterMiddleware(http.DefaultServeMux)

    http.ListenAndServe(":8080", httpHandler)

}
```

## Performance

General Load:
- **Target Rate**: 25,000 req/s
- **Duration**: 30 seconds
- **Total Requests**: 750,000

General Summary:

| Metric              | Value                             |
|---------------------|-----------------------------------|
| Request Rate        | 25,000 req/s                      |
| Actual Throughput   | 13,667 req/s (54.67% allowed)     |
| Total Requests      | 750,000                           |
| Successful (200)    | 410,000 (54.67%)                  |
| Rate Limited (429)  | 340,000 (45.33%)                  |


Latency:

| Percentile | Latency   |
|------------|-----------|
| **Mean**   | 53.4 µs   |
| **p50**    | 48.2 µs   |
| **p95**    | 73.9 µs   |
| **p99**    | 182.9 µs  |
| **Max**    | 5.1 ms    |


Normal Load:
- **Target Rate**: 10,000 req/s
- **Duration**: 30 seconds
- **Total Requests**: 300,000

General Summary:

| Metric              | Value                             |
|---------------------|-----------------------------------|
| Request Rate        | 10,000 req/s                      |
| Actual Throughput   | 10,000 req/s                      |
| Total Requests      | 30,000                            |
| Successful (200)    | 30,000 (100%)                     |

Latency:

| Percentile | Latency   |
|------------|-----------|
| **Mean**   | 45.4 µs   |
| **p50**    | 40.4 µs   |
| **p95**    | 63.2 µs   |
| **p99**    | 144.7 µs  |
| **Max**    | 3.0 ms    |
