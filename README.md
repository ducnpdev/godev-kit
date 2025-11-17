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

### 🆕 New Feature: Auto Payment System

The latest addition to GoDev Kit is a complete **Auto Payment System** for electric bills with:

- **RESTful API**: Register payments via HTTP endpoints
- **Database Persistence**: PostgreSQL for payment data storage
- **Asynchronous Processing**: Kafka-based payment processing
- **Real-time Status**: Payment status tracking and updates
- **User Management**: Payment history per user
- **Swagger Documentation**: Complete API documentation
- **Production Ready**: Error handling, logging, and monitoring

This demonstrates how to build a complete business feature following clean architecture principles with proper separation of concerns, database integration, message queuing, and comprehensive documentation.

### Key Features

- **Clean Architecture**: Separation of concerns between controllers, use cases, repositories, and entities.
- **Configurable**: Centralized configuration via YAML and environment variables.
- **Database Integration**: Built-in support for PostgreSQL, Redis, and migration scripts.
- **Kafka Integration**: Easily produce and consume messages with built-in Kafka support.
- **NATS Integration**: Built-in support for NATS messaging for event-driven architectures.
- **Redis Integration**: Use Redis for caching or fast key-value storage with ready-to-use modules.
- **User Login Module**: Includes JWT-based authentication and user management out of the box.
- **Payment System**: Complete auto payment system for electric bills with Kafka processing.
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
  - [Kafka Integration](#kafka-integration)
  - [Redis Integration](#redis-integration)
  - [User Login Module](#user-login-module)
  - [Payment System](#payment-system)
    - [Overview](#payment-overview)
    - [Architecture](#payment-architecture)
    - [API Endpoints](#payment-api-endpoints)
    - [Database Schema](#payment-database-schema)
    - [Kafka Processing](#payment-kafka-processing)
    - [Setup and Usage](#payment-setup-and-usage)
  - [VietQR Integration](#vietqr-integration)
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

### Payment System

A complete auto payment system for electric bills with RESTful API, PostgreSQL database, Kafka processing, and Swagger documentation.

**Key Features:**
- ✅ Payment registration via HTTP API
- ✅ PostgreSQL database with migrations
- ✅ Kafka asynchronous processing
- ✅ Real-time status tracking
- ✅ User payment history
- ✅ Complete Swagger documentation

**API Endpoints:**
- `POST /v1/payments` - Register a new payment
- `GET /v1/payments/{id}` - Get payment by ID  
- `GET /v1/users/{user_id}/payments` - Get user payment history

**Quick Start:**
```bash
# Run database migrations
make migrate-up

# Start application
go run cmd/app/main.go

# Access Swagger UI
http://localhost:10000/swagger/index.html
```

📖 **Detailed Documentation:**
- [PAYMENT_SYSTEM.md](PAYMENT_SYSTEM.md) - System architecture and overview
- [PAYMENT_SETUP.md](PAYMENT_SETUP.md) - Complete setup instructions
- [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) - Swagger UI usage guide
- [PAYMENT_SUMMARY.md](PAYMENT_SUMMARY.md) - Implementation summary
- [examples/payment_demo.go](examples/payment_demo.go) - Usage examples

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
```bash
# Generate or update Swagger docs (OpenAPI)
swag init -g internal/controller/http/router.go
# This will update docs/swagger.json, docs/swagger.yaml, and docs/docs.go
```

4. Access Swagger UI
If enabled in your config, you can view the interactive API docs at:
```
http://localhost:8080/swagger/index.html
```

5. Writing Annotations
- Use swaggo/swag annotations in your handler functions for automatic doc generation.
- See existing handlers in `internal/controller/http/v1/` for examples.

6. Update Documentation
- Re-run the `swag init` command after adding or changing API endpoints or annotations.

### Available Endpoints
- All available endpoints are documented in the Swagger UI and OpenAPI files.
- **New:** `GET /v1/redis/shipper/location/:shipper_id` — Get the latest location of a shipper (cache-aside pattern).
- **Payment System:** Complete payment endpoints for electric bill processing:
  - `POST /v1/payments` — Register a new payment
  - `GET /v1/payments/{id}` — Get payment by ID
  - `GET /v1/users/{user_id}/payments` — Get payments by user ID

## API Documentation

- The OpenAPI/Swagger spec is always available in `docs/swagger.json` and `docs/swagger.yaml`.
- To regenerate after code changes, run:
  ```bash
  swag init -g cmd/app/main.go
  ```
- For interactive docs, visit `/swagger/index.html` when running the service.
- **Payment System Documentation**: Complete Swagger documentation for payment endpoints is available at `http://localhost:10000/swagger/index.html`

## Development

### Prerequisites

- Go 1.20 or higher
- Docker (for local development)
- Make (for building and running)
- Git (for version control)

### Building

```bash
# Build the application
make build

# Build the Docker image
make docker-build
```

### Running

```bash
# Run the application locally
make run

# Run the application in Docker
make docker-run
```

### Testing

The project includes comprehensive unit tests for API routes and handlers. Use the following Makefile commands to run tests:

```bash
# Run all tests with race detector and coverage
make test

# Run unit tests only
make test-unit

# Run tests with coverage report (generates coverage.html)
make test-coverage

# Run tests with coverage and automatically open HTML report
make test-coverage-view

# Run tests for a specific package
# Example: make test-package PACKAGE=./internal/controller/http/v1
make test-package PACKAGE=./internal/controller/http/v1

# Run tests without race detector (faster execution)
make test-short

# Run tests with verbose output
make test-verbose

# Run integration tests
make integration-test
```

#### Test Coverage

The test suite includes:
- **API Route Tests**: Comprehensive unit tests for all HTTP routes in `internal/controller/http/v1/`
  - User routes: Create, Get, List, Update, Delete, Login (22 test cases, ~24% coverage)
- **Use Case Tests**: Unit tests for business logic in `internal/usecase/`
  - Translation: History function (7 test cases)
  - VietQR: GenerateQR and InquiryQR functions (14 test cases, 92.3% coverage)
- **Mock Dependencies**: Uses testify/mock for mocking repositories, use cases, and logger
- **Test Scenarios**: Success cases, validation errors, service errors, and edge cases
- **Coverage Reports**: Generate HTML coverage reports with `make test-coverage`

#### Running Tests with `go test`

You can also run tests directly using `go test` commands:

```bash
# Run all tests in the project
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run tests for a specific package
go test ./internal/controller/http/v1
go test ./internal/usecase/translation
go test ./internal/usecase/vietqr

# Run a specific test function
go test -v ./internal/controller/http/v1 -run TestNewUserRoutes
go test -v ./internal/usecase/translation -run TestUseCase_History
go test -v ./internal/usecase/vietqr -run TestVietQRUseCase_GenerateQR

# Run tests with coverage for specific package
go test -cover ./internal/usecase/vietqr
go test -cover ./internal/usecase/translation

# Run tests with detailed coverage report
go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/usecase/vietqr && \
go tool cover -html=coverage.txt -o coverage.html
```

#### Example: Running User Route Tests

```bash
# Test only the user routes package
make test-package PACKAGE=./internal/controller/http/v1

# Test specific test function with verbose output
go test -v ./internal/controller/http/v1 -run TestNewUserRoutes

# Test with coverage
go test -cover ./internal/controller/http/v1 -run TestNewUserRoutes

# Test with full coverage report for specific package
go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/controller/http/v1 && \
go tool cover -html=coverage.txt -o coverage.html
```

#### Example: Running Use Case Tests

```bash
# Run translation use case tests
go test -v ./internal/usecase/translation -run TestUseCase_History
go test -cover ./internal/usecase/translation

# Run VietQR use case tests
go test -v ./internal/usecase/vietqr -run TestVietQRUseCase
go test -cover ./internal/usecase/vietqr -run TestVietQRUseCase_GenerateQR
go test -cover ./internal/usecase/vietqr -run TestVietQRUseCase_InquiryQR

# Run all use case tests
go test -v ./internal/usecase/...
go test -cover ./internal/usecase/...
```

#### Pre-commit Testing

The `pre-commit` target runs all checks including tests:
```bash
make pre-commit
```

## Deployment

### Prerequisites

- A Kubernetes cluster (e.g., Minikube, Kind, EKS, GKE)
- kubectl (for Kubernetes commands)
- Helm (for deploying with Helm charts)
- Docker (for building images)

### Steps

1. **Build and Push Images**:
   ```bash
   # Build and push the application image
   make docker-build-push
   ```

2. **Deploy to Kubernetes**:
   ```bash
   # Create namespace (if not exists)
   kubectl create namespace godev-kit

   # Apply Helm chart
   helm install godev-kit ./helm/godev-kit
   ```

3. **Expose Services**:
   ```bash
   # Expose the application using Ingress or Service
   kubectl expose deployment godev-kit --type=LoadBalancer --port=80 --target-port=8080
   ```

4. **Access the Application**:
   ```bash
   # Get the external IP
   kubectl get svc godev-kit
   ```

## NATS Integration

### Running a NATS Server

```bash
# Start a NATS server locally
nats-server
```

### Configuring NATS

```yaml
NATS:
  URL: nats://localhost:4222
  CLUSTER:
    NAME: godev-kit
    PEERS:
      - nats://localhost:4223
      - nats://localhost:4224
  AUTH:
    USER: godev-kit
    PASSWORD: godev-kit
```

### NATS API Usage

- The NATS server exposes a gRPC interface for managing JetStream streams and consumers.
- You can use the `nats` CLI tool or a gRPC client to interact with it.
- Example:
  ```bash
  # List streams
  nats stream ls

  # Create a stream
  nats stream add my_stream --subjects="my.subject" --storage=file --file="nats://localhost:4222/jetstream/my_stream"
  ```

## Payment System Files

**New Files Created:** 14 files including core components, database migrations, and documentation  
**Modified Files:** 5 existing files for integration  
**Total:** ~2000+ lines of code

📁 **File Structure:** See [PAYMENT_SUMMARY.md](PAYMENT_SUMMARY.md) for complete file listing and implementation details.

## Contributing

1. Fork the repository.
2. Create a new branch for your feature.
3. Make your changes and commit them.
4. Push to your branch.
5. Create a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
