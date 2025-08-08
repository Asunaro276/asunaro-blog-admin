# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language

All outputs and responses should be provided in Japanese.


## Project Overview

CMS API is a REST API system built with Go and Echo framework, following Clean Architecture principles. It provides content management functionality through a serverless architecture using AWS Lambda + API Gateway + Aurora Serverless v2 (PostgreSQL).

## Development Commands

### Build and Run
- `make run` - Run API server locally (port :8080)
- `go run cmd/main.go` - Run standalone development server
- `go run cmd/lambda/main.go` - Run Lambda handler for testing

### Testing
- `make test` - Run all tests
- `make mock` - Generate mocks using Mockery

## Architecture

### Clean Architecture Structure
```
internal/
├── domain/entity/      # Domain entities (Content, etc.)
├── usecase/           # Business logic layer
│   ├── content/       # Content-related use cases
│   └── healthcheck/   # Health check functionality  
├── infrastructure/    # Infrastructure layer
│   ├── repository/    # Data access (Aurora Serverless v2)
│   └── controller/    # HTTP handlers
└── di/               # Dependency injection and routing
    └── route.go      # Echo router configuration
```

### Key Technologies
- **Language**: Go 1.24.1
- **Web Framework**: Echo v4.13.4
- **Database**: Aurora Serverless v2 (PostgreSQL 15)
- **Architecture**: Clean Architecture + Serverless

## API Endpoints

Based on the design specifications in `/docs/design/cms-api/`:

### Core Endpoints
- `GET /contents/{id}` - Retrieve specific content details
- `GET /contents` - Retrieve content list with pagination and filtering
- `GET /healthcheck` - System health check

### Query Parameters (for /contents)
- `limit` - Number of items (1-100, default: 20)
- `offset` - Pagination offset (default: 0)  
- `status` - Filter by status (`draft`, `published`, `archived`)
- `category` - Filter by category
- `tags` - Filter by tags (comma-separated)
- `search` - Search in title/content
- `sort` - Sort field (`createdAt`, `updatedAt`, `publishedAt`, `title`)
- `order` - Sort order (`asc`, `desc`)

### Response Format
```json
{
  "success": boolean,
  "data": object | array | null,
  "error": {
    "code": "string",
    "message": "string", 
    "details": object,
    "timestamp": "string (ISO 8601)",
    "requestId": "string"
  },
  "meta": {
    "requestId": "string",
    "timestamp": "string (ISO 8601)",
    "processingTimeMs": number
  }
}
```

## Deployment Environments

### Local Development
- Entry point: `cmd/main.go`
- Server runs on `:8080`
- Uses local Aurora connection

### AWS Lambda Production
- Entry point: `cmd/lambda/main.go`
- Deployed via API Gateway
- Uses Aurora Serverless v2 with Secrets Manager

## Testing Strategy

### Test Patterns
- **Table-driven tests** - Standard Go testing pattern for multiple test cases
- **Mock generation** - Uses Mockery for repository and external dependency mocks
- **Test containers** - Integration testing with real database instances

### Test Structure Example
```go
tests := []struct {
    name           string
    input          string
    expectedOutput string
}{
    {
        name:           "Valid case",
        input:          "test-input",
        expectedOutput: "expected-output",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

## Non-Functional Requirements

Based on `/docs/spec/cms-api-requirements.md`:

### Performance
- 95% of requests respond within 1 second
- Support 1000 concurrent connections
- Database queries complete within 500ms

### Security
- HTTPS communication only
- SQL injection prevention via parameterized queries
- No sensitive information in logs
- Proper CORS configuration

### Availability
- 99.9% uptime requirement
- No single points of failure
- Automatic recovery capabilities

## Error Handling

### Standard Error Codes
- `INVALID_PARAMETER` (400) - Invalid parameter values
- `CONTENT_NOT_FOUND` (404) - Content not found
- `INTERNAL_ERROR` (500) - Internal server error
- `DATABASE_ERROR` (500) - Database connectivity issues
- `SERVICE_UNAVAILABLE` (503) - Service temporarily unavailable

## Development Notes

### Database Migration
- **Current**: Transitioning from DynamoDB to Aurora Serverless v2
- **Target**: PostgreSQL 15 with proper relational schema
- Repository layer abstracted to support both during migration

### Code Conventions
- Follow standard Go conventions and gofmt formatting
- Use meaningful variable and function names
- Implement comprehensive error handling
- Write tests for all business logic
- Use structured logging (JSON format)

### Dependencies Management
- Use Go modules for dependency management
- Keep dependencies up to date
- Prefer standard library when possible
- Use Mockery for generating test mocks
