# API Naming Conventions

## General Principles

1. **Consistency**: Use consistent naming across all APIs
2. **Clarity**: Names should be self-explanatory
3. **REST Standards**: Follow RESTful conventions
4. **Versioning**: Use API versioning for compatibility

## URL Structure

### Base URL Pattern
```
https://api.domain.com/api/v{version}/{resource}
```

### Examples
```
https://api.domain.com/api/v1/users
https://api.domain.com/api/v1/products
https://api.domain.com/api/v1/orders
```

## HTTP Methods

### Standard CRUD Operations
- **GET**: Retrieve resource(s)
- **POST**: Create new resource
- **PUT**: Update entire resource
- **PATCH**: Partially update resource
- **DELETE**: Remove resource

### Method Usage Examples
```
GET    /api/v1/users          # List all users
GET    /api/v1/users/{id}     # Get specific user
POST   /api/v1/users          # Create new user
PUT    /api/v1/users/{id}     # Update user (full replacement)
PATCH  /api/v1/users/{id}     # Partial update user
DELETE /api/v1/users/{id}     # Delete user
```

## Resource Naming

### Resource Names
- Use **plural nouns** for collections: `/users`, `/products`, `/orders`
- Use **lowercase** with hyphens for multi-word resources: `/user-profiles`, `/order-items`
- Use **nouns**, not verbs: `/users` not `/get-users`

### Examples
```
✅ Good:
/api/v1/users
/api/v1/products
/api/v1/order-items
/api/v1/user-profiles

❌ Bad:
/api/v1/user          # singular
/api/v1/getUsers      # verb
/api/v1/user_profiles # underscore
/api/v1/Users         # uppercase
```

## Query Parameters

### Filtering
- Use descriptive parameter names
- Support multiple filters
- Use consistent filter patterns

```
GET /api/v1/users?status=active
GET /api/v1/products?category=electronics&price_min=100&price_max=500
GET /api/v1/orders?created_after=2023-01-01&status=pending
```

### Pagination
- Use `page` and `limit` parameters
- Provide pagination metadata in response

```
GET /api/v1/users?page=1&limit=20
GET /api/v1/products?page=2&limit=50
```

### Sorting
- Use `sort` parameter with field names
- Support ascending/descending order

```
GET /api/v1/users?sort=created_at
GET /api/v1/users?sort=-created_at  # descending
GET /api/v1/users?sort=name,created_at
```

### Searching
- Use `q` or `search` parameter for general search
- Use specific field names for targeted search

```
GET /api/v1/users?q=john
GET /api/v1/users?name=john&email=john@example.com
```

## Path Parameters

### Resource Identifiers
- Use meaningful parameter names
- Use consistent identifier patterns

```
GET /api/v1/users/{user_id}
GET /api/v1/users/{user_id}/orders
GET /api/v1/users/{user_id}/orders/{order_id}
```

### Nested Resources
- Limit nesting to 2-3 levels maximum
- Use clear hierarchy

```
GET /api/v1/users/{user_id}/orders
GET /api/v1/users/{user_id}/orders/{order_id}/items
```

## Action Endpoints

### Non-CRUD Operations
- Use verbs for actions that don't fit CRUD
- Place actions after resource identifier

```
POST /api/v1/users/{user_id}/activate
POST /api/v1/users/{user_id}/deactivate
POST /api/v1/orders/{order_id}/cancel
POST /api/v1/orders/{order_id}/ship
```

### Batch Operations
- Use plural forms for batch operations
- Use descriptive action names

```
POST /api/v1/users/bulk-create
POST /api/v1/users/bulk-update
POST /api/v1/users/bulk-delete
```

## HTTP Status Codes

### Success Responses
- **200 OK**: Successful GET, PUT, PATCH
- **201 Created**: Successful POST (resource created)
- **202 Accepted**: Async operation accepted
- **204 No Content**: Successful DELETE

### Client Error Responses
- **400 Bad Request**: Invalid request format
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Access denied
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict
- **422 Unprocessable Entity**: Validation error

### Server Error Responses
- **500 Internal Server Error**: Server error
- **502 Bad Gateway**: Upstream error
- **503 Service Unavailable**: Service temporarily unavailable

## Request/Response Format

### Request Headers
```
Content-Type: application/json
Accept: application/json
Authorization: Bearer {token}
```

### Response Headers
```
Content-Type: application/json
X-Request-ID: {request_id}
X-Rate-Limit-Remaining: {count}
```

### JSON Field Naming
- Use **snake_case** for JSON field names
- Use consistent naming across all endpoints

```json
{
  "user_id": 123,
  "first_name": "John",
  "last_name": "Doe",
  "email_address": "john@example.com",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

## Error Response Format

### Standard Error Structure
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email_address",
        "code": "INVALID_FORMAT",
        "message": "Invalid email format"
      }
    ],
    "request_id": "req_123456789"
  }
}
```

### Error Codes
- Use **UPPER_SNAKE_CASE** for error codes
- Use descriptive, consistent error codes

```
VALIDATION_ERROR
RESOURCE_NOT_FOUND
UNAUTHORIZED_ACCESS
INSUFFICIENT_PERMISSIONS
RATE_LIMIT_EXCEEDED
INTERNAL_SERVER_ERROR
```

## Authentication and Authorization

### Authentication Header
```
Authorization: Bearer {jwt_token}
```

### Protected Endpoints
- Clearly indicate which endpoints require authentication
- Use consistent authentication mechanism

```
GET /api/v1/users/profile      # Requires authentication
POST /api/v1/users/logout      # Requires authentication
POST /api/v1/auth/login        # Public endpoint
```

## Rate Limiting

### Rate Limit Headers
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```

### Rate Limit Response
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again later.",
    "retry_after": 60
  }
}
```

## Versioning

### URL Versioning
- Use semantic versioning: v1, v2, v3
- Maintain backward compatibility within major versions

```
/api/v1/users
/api/v2/users
```

### Version Deprecation
- Provide deprecation warnings
- Give sufficient notice before removal

```
Deprecation: version="v1", date="2024-01-01", link="/api/v2/users"
```

## Documentation

### OpenAPI/Swagger
- Document all endpoints
- Include request/response examples
- Specify required/optional parameters

### Example Documentation
```yaml
paths:
  /api/v1/users:
    get:
      summary: List users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: object
                properties:
                  users:
                    type: array
                    items:
                      $ref: '#/components/schemas/User'
                  pagination:
                    $ref: '#/components/schemas/Pagination'
```

## Testing

### Endpoint Testing
- Test all HTTP methods
- Test success and error scenarios
- Test edge cases and validation

### Test Examples
```go
func TestGetUser_Success(t *testing.T) {
    // Test successful user retrieval
}

func TestGetUser_NotFound(t *testing.T) {
    // Test user not found scenario
}

func TestCreateUser_ValidationError(t *testing.T) {
    // Test validation error scenarios
}
```