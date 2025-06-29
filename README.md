![Go Dev Kit Template](docs/img/godevkit-logo.svg)
- copy from: go-clean-template
# Go Dev Kit template

Godev Kit template for Golang services

## Contact
If you implement a new feature or need support running the project, please contact me!
- Email: ducnp09081998@gmail.com
- FB: https://www.facebook.com/phucducdev
- Linkedin: https://www.linkedin.com/in/phucducktpm/
- Or open an issue or pull request on GitHub.

## Overview

**GoDev Kit** is a modular, production-ready template for building robust Golang microservices and backend applications. It provides a clean architecture foundation, best practices, and ready-to-use integrations for common infrastructure components, allowing you to focus on your business logic instead of boilerplate setup.

### Key Features

- **Clean Architecture**: Separation of concerns between controllers, use cases, repositories, and entities.
- **Configurable**: Centralized configuration via YAML and environment variables.
- **Database Integration**: Built-in support for PostgreSQL, Redis, and migration scripts.
- **Messaging**: Ready-to-use Kafka, RabbitMQ, and NATS integrations for event-driven architectures.
- **API Ready**: HTTP and gRPC server templates, with Swagger/OpenAPI documentation.
- **Observability**: Prometheus metrics and structured logging out of the box.
- **Extensible**: Easily add new features, endpoints, or infrastructure components.
- **Developer Experience**: Makefile tasks, Docker support, and example code for rapid development.

### Use Cases

- Rapidly bootstrap new Go microservices or backend APIs.
- Learn and apply best practices for scalable Go service design.
- Serve as a reference implementation for clean, maintainable Go codebases.

## Table of Contents
- [Quick Start](#quick-start)
- [Features](#features)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [YAML Configuration](#yaml-configuration)
  - [Prometheus Metrics](#prometheus-metrics)
    - [Enable/Disable Metrics](#enable-disable-metrics)
    - [Configure Skip Paths](#configure-skip-paths)
    - [Example Output](#example-output)
  - [Swagger Documentation](#swagger)
    - [Installation](#installation)
    - [Generate Documentation](#generate-documentation)
    - [Access Swagger UI](#access-swagger-ui)
    - [Writing Annotations](#writing-swagger-annotations)
    - [Available Endpoints](#available-endpoints)
    - [Update Documentation](#update-documentation)
    - [Swagger UI Features](#swagger-ui-features)
- [Project Structure](#project-structure)
- [API Documentation](#api-documentation)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Building](#building)
  - [Running](#running)
  - [Testing](#testing)
- [Deployment](#deployment)
- [NATS Integration](#nats-integration)
  - [Running a NATS Server](#running-a-nats-server)
  - [Configuring NATS](#configuring-nats)
  - [NATS API Usage](#nats-api-usage)
- [Contributing](#contributing)
- [License](#license)

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

2. Installation
```bash
# Install Swagger CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
swag --version
```

3. Generate Documentation
```