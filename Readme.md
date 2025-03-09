# Golang REST API with Dependency Injection

This is a boilerplate for a Go REST API with dependency injection pattern, following clean architecture principles.

## Features

- Clean architecture with separation of concerns
- Dependency injection for better testability
- PostgreSQL database with connection pooling
- Structured logging with zerolog
- Graceful shutdown
- Configuration using Viper (file + environment variables)
- HTTP server using gorilla/mux router
- Complete user CRUD functionality
- Health check endpoint
- Docker support
- Database migrations

## Project Structure

```
├── cmd
│   └── api
│       └── main.go
├── config
│   └── config.go
├── internal
│   ├── app
│   │   └── app.go
│   ├── handler
│   │   ├── handlers.go
│   │   ├── health.go
│   │   └── user.go
│   ├── model
│   │   └── user.go
│   ├── repository
│   │   ├── health.go
│   │   ├── repositories.go
│   │   └── user.go
│   └── service
│       ├── health.go
│       ├── services.go
│       └── user.go
├── migrations
│   ├── 000001_create_users_table.down.sql
│   └── 000001_create_users_table.up.sql
├── pkg
│   ├── database
│   │   └── database.go
│   ├── httputil
│   │   └── response.go
│   └── logger
│       └── logger.go
├── config.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate) (for running migrations)

### Setup

1. Clone the repository:
```shell
git clone https://github.com/yourusername/myapi.git
cd myapi
```

2. Install dependencies:
```shell
go mod download
```

3. Create PostgreSQL database:
```shell
createdb myapp
```

4. Run migrations:
```shell
make migration-up
```

5. Run the application:
```shell
make run
```

The API will be available at `http://localhost:8080`.

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Docker

Build and run using Docker:

```shell
make docker-build
make docker-run
```

## Configuration

The application can be configured using the `config.yaml` file or environment variables.

Environment variables:
- `SERVER_ADDRESS` - Server address (default: `:8080`)
- `DATABASE_URL` - Database connection string
- `LOGGER_LEVEL` - Log level (debug, info, warn, error, fatal)

## License

This project is licensed under the MIT License - see the LICENSE file for details.