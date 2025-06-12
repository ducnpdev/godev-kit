![Go Dev Kit Template](docs/img/godevkit-logo.svg)
- copy from: go-clean-template
# Go Dev Kit template

Godev Kit template for Golang services

## Overview

todo

## Content

todo

## Quick start

## Feature
### Config data: 
- handle load env from yaml file.
  - config struct into file `config/config.go`
  - value yaml into file `config/config.yaml`

### Prometheus Metrics:
1. On|off metrics
```yaml
METRICS:
  ENABLED: true|false
```
2. Config bypass route api.
```yaml
METRICS:
  ...
  SKIP_PATHS: "/swagger/*;/metrics"
```
- Remove some paths from metrics with sep ";"
```go
prometheus.SetSkipPaths(strings.Split(cfg.Metrics.SetSkipPaths, ";"))
```
3. Code example: https://github.com/ansrivas/fiberprometheus
4. Output: 
- access browser: http://127.0.0.1:8080/metrics
```text
# HELP http_request_duration_seconds Duration of all HTTP requests by status code, method and path.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/metrics",service="godev-kit",status_code="200",le="1e-09"} 0
http_request_duration_seconds_bucket{method="GET",path="/metrics",service="godev-kit",status_code="200",le="60"} 1
http_request_duration_seconds_bucket{method="GET",path="/metrics",service="godev-kit",status_code="200",le="+Inf"} 1
http_request_duration_seconds_sum{method="GET",path="/metrics",service="godev-kit",status_code="200"} 0.000803375
http_request_duration_seconds_count{method="GET",path="/metrics",service="godev-kit",status_code="200"} 1
# HELP http_requests_in_progress_total All the requests in progress
# TYPE http_requests_in_progress_total gauge
http_requests_in_progress_total{method="GET",service="godev-kit"} 1
# HELP http_requests_total Count all http requests by status code, method and path.
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/metrics",service="godev-kit",status_code="200"} 1
```

### Swagger
1. On|Off
```yaml
SWAGGER:
  ENABLED: true|false
```