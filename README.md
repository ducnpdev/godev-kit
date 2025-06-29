![Go Dev Kit Template](docs/img/godevkit-logo.svg)
- copy from: go-clean-template
# Go Dev Kit template

Godev Kit template for Golang services

## Overview

todo

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
```bash
# Generate Swagger docs
swag init -g cmd/app/main.go -o docs

# This will create:
# - docs/docs.go
# - docs/swagger.json
# - docs/swagger.yaml
```

4. Access Swagger UI
- Start your application:
```bash
go run cmd/app/main.go
```
- Open your browser and navigate to: `http://localhost:8080/swagger/`

5. Writing Swagger Annotations
```go
// @Summary     Create user
// @Description Create a new user
// @ID          create-user
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Param       request body request.CreateUser true "Create user"
// @Success     201 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user [post]
func (r *V1) createUser(ctx *fiber.Ctx) error {
    // ... handler implementation
}
```

6. Available Endpoints
- User Management:
  - POST /v1/user - Create user
  - GET /v1/user - List users
  - GET /v1/user/{id} - Get user by ID
  - PUT /v1/user/{id} - Update user
  - DELETE /v1/user/{id} - Delete user
- Translation Service:
  - POST /v1/translation/do-translate - Translate text
  - GET /v1/translation/history - Show translation history

7. Update Documentation
- After making changes to your API endpoints or models, regenerate the Swagger docs:
```bash
swag init -g cmd/app/main.go -o docs
```

8. Swagger UI Features
- Interactive API documentation
- Try out API endpoints directly from the browser
- View request/response schemas
- Download OpenAPI specification (JSON/YAML)
```

## Project Structure
```
.
├── cmd/                    # Application entry points
│   └── app/               # Main application
│       ├── config/        # Configuration
│       └── main.go        # Application entry point
├── config/                # Configuration files
│   ├── config.go         # Configuration structure
│   └── config.yaml       # Configuration values
├── docs/                  # Documentation
│   ├── img/              # Images
│   ├── docs.go           # Swagger documentation
│   ├── swagger.json      # OpenAPI JSON
│   └── swagger.yaml      # OpenAPI YAML
├── internal/             # Private application code
│   ├── controller/       # API handlers
│   ├── entity/          # Business entities
│   ├── repo/            # Repository layer
│   └── usecase/         # Business logic
├── migrations/           # Database migrations
├── pkg/                  # Public library code
├── vendor/              # Application dependencies
├── .github/             # GitHub templates and workflows
├── .vscode/             # VS Code settings
├── nginx/               # Nginx configuration
├── .dockerignore        # Docker ignore file
├── .gitignore          # Git ignore file
├── .golangci.yml       # Golang linter config
├── go.mod              # Go module file
├── go.sum              # Go module checksum
├── LICENSE             # License file
├── Makefile            # Build automation
└── README.md           # Project documentation
```

## API Documentation
The API documentation is available through Swagger UI when the application is running. See the [Swagger Documentation](#swagger) section for details.

## Development

### Prerequisites
- Go 1.21 or higher
- PostgreSQL
- RabbitMQ (for RPC)
- NATS (for messaging, can be run with Docker)
- Make (optional, for using Makefile)

### Building
```bash
# Build the application
make build

# Or manually
go build -o bin/app cmd/app/main.go
```

### Running
```bash
# Run the application
make run

# Or manually
go run cmd/app/main.go
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Deployment
The application can be deployed using Docker:

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

## NATS Integration

### Running a NATS Server
You can quickly start a local NATS server using Docker:

```bash
make docker-run-nats
```

This will run the official NATS server on port 4222.

To stop and remove the container:
```bash
docker stop nats-server && docker rm nats-server
```

### Configuring NATS
NATS connection settings are managed in `config/config.yaml`:

```yaml
NATS:
  URL: nats://localhost:4222
  TIMEOUT: 3s
```

### NATS API Usage
The service exposes HTTP endpoints for publishing and subscribing to NATS subjects:

- **Publish message:**
  - `POST /v1/nats/publish/{subject}`
  - Body: `{ "data": "your message" }`

- **Subscribe to subject:**
  - `GET /v1/nats/subscribe/{subject}`
  - Returns the first message received on the subject (demo purpose)

You can try these endpoints via Swagger UI at `http://localhost:8080/swagger/` when the app is running.

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

If you implement a new feature or need support running the project, please contact me!

- Email: ducnp09081998@gmail.com
- FB: https://www.facebook.com/phucducdev
- Linkedin: https://www.linkedin.com/in/phucducktpm/
- Or open an issue or pull request on GitHub.