# Coding Standards

## General Principles

1. **Clean Code**: Write self-documenting, readable code
2. **SOLID Principles**: Follow Single Responsibility, Open/Closed, Liskov Substitution, Interface Segregation, and Dependency Inversion
3. **DRY (Don't Repeat Yourself)**: Avoid code duplication
4. **KISS (Keep It Simple, Stupid)**: Write simple, straightforward solutions
5. **YAGNI (You Aren't Gonna Need It)**: Don't implement functionality until it's needed

## Code Organization

### Package Structure
```
internal/
├── app/            # Application initialization and setup
├── controller/     # HTTP/gRPC controllers (adapters)
├── entity/         # Business entities and domain models
├── repo/           # Data access layer (repositories)
└── usecase/        # Business logic (use cases)

pkg/
├── httpserver/     # HTTP server utilities
├── grpcserver/     # gRPC server utilities
├── postgres/       # PostgreSQL utilities
├── redis/          # Redis utilities
├── kafka/          # Kafka utilities
├── nats/           # NATS utilities
├── rabbitmq/       # RabbitMQ utilities
└── logger/         # Logging utilities
```

### Layer Dependencies
- **Controllers** depend on **Use Cases**
- **Use Cases** depend on **Entities** and **Repository interfaces**
- **Repositories** depend on **Entities**
- **Entities** have no dependencies

## Naming Conventions

### Files and Directories
- Use snake_case for file names: `user_repository.go`
- Use lowercase for directory names: `usecase`, `controller`
- Use descriptive names: `translation_service.go` not `service.go`

### Variables and Functions
- Use camelCase for variables and functions: `userService`, `getUserByID`
- Use PascalCase for exported functions: `NewUserService`, `GetUserByID`
- Use meaningful names: `userID` not `id`, `userRepository` not `repo`

### Constants
- Use UPPER_SNAKE_CASE: `MAX_RETRY_ATTEMPTS`, `DEFAULT_TIMEOUT`
- Group related constants in const blocks

### Interfaces
- Use descriptive names ending with interface behavior: `UserRepository`, `TranslationService`
- Keep interfaces small and focused (Interface Segregation Principle)

## Code Style

### Formatting
- Use `gofmt` or `goimports` to format code
- Line length should not exceed 120 characters
- Use tabs for indentation, not spaces

### Comments
- Use // for single-line comments
- Use /* */ for multi-line comments
- Document all exported functions, types, and packages
- Use godoc format for documentation

### Error Handling
- Always handle errors explicitly
- Use meaningful error messages
- Wrap errors with context using `fmt.Errorf` or `errors.Wrap`
- Don't ignore errors with `_`

### Function Design
- Keep functions small and focused (Single Responsibility Principle)
- Maximum 20 lines per function (excluding comments)
- Use early returns to avoid deep nesting
- Limit function parameters to 3-4 maximum

## Testing

### Test Files
- Name test files with `_test.go` suffix
- Place test files in the same package as the code being tested
- Use descriptive test function names: `TestUserService_GetUserByID_Success`

### Test Structure
- Use Table-driven tests for multiple test cases
- Follow AAA pattern: Arrange, Act, Assert
- Use testify/assert for assertions
- Mock external dependencies

### Coverage
- Aim for 80%+ code coverage
- Focus on testing business logic (use cases)
- Test both success and failure scenarios

## Database and Repository Pattern

### Repository Interface
```go
type UserRepository interface {
    Create(ctx context.Context, user entity.User) (entity.User, error)
    GetByID(ctx context.Context, id int64) (entity.User, error)
    Update(ctx context.Context, user entity.User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter entity.UserFilter) ([]entity.User, error)
}
```

### SQL Queries
- Use parameterized queries to prevent SQL injection
- Use Squirrel query builder for complex queries
- Keep SQL queries in repository implementations
- Use database transactions for multi-step operations

## HTTP API Design

### REST Conventions
- Use HTTP verbs correctly: GET, POST, PUT, DELETE
- Use plural nouns for resource names: `/users`, `/products`
- Use consistent URL patterns: `/api/v1/users/{id}`
- Return appropriate HTTP status codes

### Request/Response
- Use JSON for request and response bodies
- Validate all input data
- Return consistent error response format
- Use pagination for list endpoints

### Middleware
- Use middleware for cross-cutting concerns: authentication, logging, metrics
- Keep middleware focused and composable
- Handle errors gracefully in middleware

## Configuration

### Environment Variables
- Use uppercase with underscores: `DATABASE_URL`, `REDIS_URL`
- Provide default values where appropriate
- Document all configuration options

### Configuration Files
- Use YAML for configuration files
- Group related configuration together
- Use environment variable overrides

## Logging

### Log Levels
- **DEBUG**: Detailed information for debugging
- **INFO**: General information about application flow
- **WARN**: Warning conditions that don't prevent operation
- **ERROR**: Error conditions that require attention
- **FATAL**: Fatal errors that cause application termination

### Log Format
- Use structured logging (JSON format)
- Include relevant context: userID, requestID, etc.
- Don't log sensitive information (passwords, tokens)
- Use consistent log message format

## Security

### Authentication & Authorization
- Use JWT tokens for authentication
- Implement proper token validation
- Use RBAC (Role-Based Access Control) for authorization
- Validate all user inputs

### Data Protection
- Never log sensitive data
- Use HTTPS for all external communications
- Encrypt sensitive data at rest
- Use secure random number generation

## Performance

### Database
- Use connection pooling
- Implement proper indexing
- Use prepared statements
- Monitor query performance

### Caching
- Use Redis for frequently accessed data
- Implement cache-aside pattern
- Set appropriate TTL values
- Handle cache invalidation properly

### Monitoring
- Implement health checks
- Use Prometheus metrics
- Monitor key performance indicators
- Set up alerting for critical issues

## Git Workflow

### Commit Messages
- Use conventional commit format: `feat: add user authentication`
- Keep commit messages concise but descriptive
- Use imperative mood: "Add feature" not "Added feature"

### Branch Naming
- Use descriptive branch names: `feature/user-authentication`
- Use prefixes: `feature/`, `bugfix/`, `hotfix/`
- Keep branch names short and focused

### Pull Requests
- Write clear PR descriptions
- Include testing information
- Request appropriate reviewers
- Ensure CI/CD passes before merging

## Documentation

### Code Documentation
- Document all exported functions and types
- Use godoc format for documentation
- Include examples in documentation
- Keep documentation up to date

### API Documentation
- Use OpenAPI/Swagger for API documentation
- Include request/response examples
- Document error responses
- Keep API docs synchronized with code

### Architecture Documentation
- Document system architecture
- Include sequence diagrams for complex flows
- Document database schema
- Keep architecture docs current