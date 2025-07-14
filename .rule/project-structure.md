# Project Structure Guidelines

## Overview

This document outlines the recommended project structure, following Clean Architecture principles and Go best practices.

## Directory Structure

```
godev-kit/
├── cmd/                    # Application entry points
│   └── app/
│       └── main.go
├── config/                 # Configuration management
│   ├── config.go
│   └── config.yaml
├── internal/              # Private application code
│   ├── app/               # Application initialization
│   ├── controller/        # Interface adapters
│   ├── entity/            # Business entities
│   ├── repo/              # Data access layer
│   └── usecase/           # Business logic
├── pkg/                   # Public library code
│   ├── httpserver/
│   ├── grpcserver/
│   ├── postgres/
│   ├── redis/
│   ├── kafka/
│   ├── nats/
│   ├── rabbitmq/
│   └── logger/
├── migrations/            # Database migrations
├── scripts/               # Build and deployment scripts
├── tests/                 # Integration and end-to-end tests
├── docs/                  # Documentation
├── .rule/                 # Project rules and guidelines
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

## Layer Definitions

### 1. Entities Layer (`internal/entity/`)

**Purpose**: Contains business objects and core business logic.

**Rules**:
- No dependencies on other layers
- Contains pure business logic
- Defines data structures and business rules

**Structure**:
```
internal/entity/
├── user.go              # User entity and business rules
├── product.go           # Product entity and business rules
├── order.go             # Order entity and business rules
├── events.go            # Domain events
├── errors.go            # Domain-specific errors
└── value_objects.go     # Value objects
```

**Example**:
```go
// user.go
package entity

type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`
    Status   UserStatus `json:"status"`
}

type UserStatus string

const (
    UserStatusActive   UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
)

func (u *User) Validate() error {
    // Business validation logic
}
```

### 2. Use Cases Layer (`internal/usecase/`)

**Purpose**: Contains application-specific business rules and orchestrates data flow.

**Rules**:
- Depends only on entities and repository interfaces
- Contains application-specific business logic
- Orchestrates data flow between entities and repositories

**Structure**:
```
internal/usecase/
├── contracts.go         # Interface definitions
├── user/
│   ├── user.go         # User use case implementation
│   └── user_test.go    # Use case tests
├── product/
│   ├── product.go
│   └── product_test.go
├── auth/
│   ├── auth.go
│   └── auth_test.go
├── kafka.go            # Kafka use case
├── redis.go            # Redis use case
└── nats.go             # NATS use case
```

**Example**:
```go
// user/user.go
package user

type UseCase struct {
    userRepo     repo.UserRepository
    logger       logger.Interface
    eventPublisher events.Publisher
}

func New(userRepo repo.UserRepository, logger logger.Interface) *UseCase {
    return &UseCase{
        userRepo: userRepo,
        logger:   logger,
    }
}

func (uc *UseCase) CreateUser(ctx context.Context, req entity.User) (entity.User, error) {
    // Business logic for creating user
}
```

### 3. Repository Layer (`internal/repo/`)

**Purpose**: Provides data access abstractions and external service integrations.

**Rules**:
- Implements repository interfaces defined in use cases
- Handles data persistence and external API calls
- Converts between domain entities and data models

**Structure**:
```
internal/repo/
├── contracts.go         # Repository interfaces
├── persistent/          # Database implementations
│   ├── user_postgres.go
│   ├── product_postgres.go
│   ├── redis.go
│   └── models/         # Database models
│       ├── user.go
│       └── product.go
├── externalapi/        # External service clients
│   ├── payment/
│   │   └── stripe.go
│   └── notification/
│       └── email.go
├── kafka.go            # Kafka repository
└── nats.go             # NATS repository
```

**Example**:
```go
// persistent/user_postgres.go
package persistent

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
    // Database implementation
}
```

### 4. Controller Layer (`internal/controller/`)

**Purpose**: Handles external interfaces (HTTP, gRPC, message queues).

**Rules**:
- Depends on use cases
- Handles request/response transformation
- Manages protocol-specific concerns

**Structure**:
```
internal/controller/
├── http/               # HTTP REST controllers
│   ├── router.go       # HTTP router setup
│   ├── middleware/     # HTTP middleware
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── recovery.go
│   └── v1/             # API version 1
│       ├── controller.go
│       ├── router.go
│       ├── user.go
│       ├── product.go
│       ├── health.go
│       ├── request/    # Request DTOs
│       │   ├── user.go
│       │   └── product.go
│       └── response/   # Response DTOs
│           ├── user.go
│           └── product.go
├── grpc/               # gRPC controllers
│   ├── router.go
│   └── v1/
│       ├── controller.go
│       ├── router.go
│       ├── user.go
│       └── response/
│           └── user.go
└── amqp_rpc/           # RabbitMQ RPC controllers
    ├── router.go
    └── v1/
        ├── controller.go
        ├── router.go
        └── user.go
```

**Example**:
```go
// http/v1/user.go
package v1

func (h *Handler) CreateUser(c *gin.Context) {
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.errorResponse(c, http.StatusBadRequest, err)
        return
    }
    
    user, err := h.userUseCase.CreateUser(c.Request.Context(), req.ToEntity())
    if err != nil {
        h.errorResponse(c, http.StatusInternalServerError, err)
        return
    }
    
    c.JSON(http.StatusCreated, response.NewUserResponse(user))
}
```

### 5. Application Layer (`internal/app/`)

**Purpose**: Wires up dependencies and initializes the application.

**Rules**:
- Handles dependency injection
- Initializes all components
- Manages application lifecycle

**Structure**:
```
internal/app/
├── app.go              # Main application setup
├── migrate.go          # Database migrations
├── wire.go             # Dependency injection (optional)
└── config.go           # Application configuration
```

**Example**:
```go
// app.go
package app

func Run(cfg *config.Config) {
    // Initialize database
    db := postgres.New(cfg.PG.URL)
    
    // Initialize repositories
    userRepo := persistent.NewUserRepository(db)
    
    // Initialize use cases
    userUseCase := user.New(userRepo, logger)
    
    // Initialize controllers
    handler := v1.NewHandler(userUseCase, logger)
    
    // Start servers
    httpServer := httpserver.New(handler, cfg.HTTP.Port)
    httpServer.Start()
}
```

## Package Organization (`pkg/`)

### Shared Libraries

**Purpose**: Reusable components that can be used across projects.

**Structure**:
```
pkg/
├── httpserver/         # HTTP server utilities
│   ├── server.go
│   ├── options.go
│   └── middleware.go
├── grpcserver/         # gRPC server utilities
│   ├── server.go
│   └── options.go
├── postgres/           # PostgreSQL utilities
│   ├── postgres.go
│   └── options.go
├── redis/              # Redis utilities
│   ├── redis.go
│   └── options.go
├── kafka/              # Kafka utilities
│   ├── consumer.go
│   ├── producer.go
│   └── manager.go
├── nats/               # NATS utilities
│   └── client.go
├── rabbitmq/           # RabbitMQ utilities
│   └── rmq_rpc/
│       ├── client/
│       └── server/
└── logger/             # Logging utilities
    └── logger.go
```

## Configuration (`config/`)

**Purpose**: Centralized configuration management.

**Structure**:
```
config/
├── config.go           # Configuration struct and loading logic
├── config.yaml         # Default configuration
├── config.local.yaml   # Local development overrides
└── config.test.yaml    # Test configuration
```

## Database Migrations (`migrations/`)

**Purpose**: Database schema versioning and migrations.

**Structure**:
```
migrations/
├── 20240101000001_create_users_table.up.sql
├── 20240101000001_create_users_table.down.sql
├── 20240101000002_create_products_table.up.sql
├── 20240101000002_create_products_table.down.sql
└── 20240101000003_add_user_indexes.up.sql
```

## Tests (`tests/`)

**Purpose**: Integration and end-to-end tests.

**Structure**:
```
tests/
├── integration/        # Integration tests
│   ├── user_test.go
│   └── product_test.go
├── e2e/               # End-to-end tests
│   ├── api_test.go
│   └── workflow_test.go
├── fixtures/          # Test data
│   ├── users.json
│   └── products.json
└── helpers/           # Test utilities
    ├── database.go
    └── server.go
```

## Documentation (`docs/`)

**Purpose**: Project documentation and API specifications.

**Structure**:
```
docs/
├── api/               # API documentation
│   ├── openapi.yaml
│   └── postman.json
├── architecture/      # Architecture diagrams
│   ├── overview.md
│   └── sequence.md
├── deployment/        # Deployment guides
│   ├── docker.md
│   └── kubernetes.md
└── development/       # Development guides
    ├── setup.md
    └── testing.md
```

## Scripts (`scripts/`)

**Purpose**: Build, deployment, and utility scripts.

**Structure**:
```
scripts/
├── build.sh           # Build scripts
├── deploy.sh          # Deployment scripts
├── test.sh            # Test scripts
├── migrate.sh         # Migration scripts
└── dev/               # Development scripts
    ├── setup.sh
    └── start.sh
```

## Rules and Guidelines (`.rule/`)

**Purpose**: Project-specific rules and coding standards.

**Structure**:
```
.rule/
├── coding-standards.md
├── api-naming-conventions.md
├── project-structure.md
└── git-workflow.md
```

## Best Practices

### 1. Dependency Direction
- Dependencies should point inward (toward entities)
- Use interfaces to invert dependencies
- Avoid circular dependencies

### 2. Error Handling
- Use domain-specific errors in entities
- Wrap errors with context in each layer
- Handle errors gracefully in controllers

### 3. Testing
- Unit tests in the same package as the code
- Integration tests in separate test packages
- Mock external dependencies

### 4. Naming Conventions
- Use descriptive package names
- Follow Go naming conventions
- Use consistent naming across layers

### 5. File Organization
- Keep related functionality together
- Separate concerns into different files
- Use subdirectories for complex domains

This structure provides a solid foundation for building scalable, maintainable Go applications while following Clean Architecture principles.