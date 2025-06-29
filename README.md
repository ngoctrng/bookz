![coverage](docs/coverage.svg) ![coverage](docs/time.svg)

# FCC Book Trading Club

The FCC Book Trading Club is a backend system designed to manage a community-driven book trading platform. Users can register, browse available books, add books to their collection, request trades with other users, and manage incoming and outgoing trade requests. The system ensures secure user authentication, maintains accurate book ownership records, and provides a seamless experience for trading books within the club.

**Project Link:** [freeCodeCamp: Manage a Book Trading Club](https://www.freecodecamp.org/learn/coding-interview-prep/take-home-projects/manage-a-book-trading-club)

## Project Structure

<pre>
.
├── cmd/                   # Application entry points
│   ├── httpserver/        # HTTP server executable
│   ├── worker/            # Background worker executable
│   └── migrate/           # Database migration tool
├── docs/                  # Documentation and OpenAPI specs
│   └── architecture.png   # Architecture diagram
│   └── diagrams/          # C4 diagrams (system-context, container, component)
├── internal/              # Private application code
│   ├── account/           # User account domain module
│   ├── book/              # Book domain module
│   ├── exchange/          # Exchange/trade domain module
│   └── ...                # Other modules
│       ├── delivery/      # HTTP handlers
│       ├── repository/    # Database implementations
│       └── usecases/      # Application logic
│       domain-model.go            # Domain Business logic
├── pkg/                   # Public shared packages
│   ├── config/            # Configuration handling
│   ├── migration/         # Database migration utilities
│   ├── hasher/            # Password hashing utilities
│   ├── token/             # JWT token utilities
│   └── testutil/          # Testing utilities
└── tools/                 # Scripts and tools
    └── compose/           # Docker compose files
</pre>

## Architecture

This project follows the Onion/Clean Architecture pattern.

![](docs/architecture.png)

Key principles:
- Dependencies flow inward
- Inner layers contain business logic
- Outer layers contain implementation details
- Domain entities are at the core
- Each layer is isolated and testable

## Documentation

For detailed architecture, system context, container, and component diagrams, as well as further technical documentation, please refer to the [docs](docs/) folder:

- [System Context Diagram](docs/diagrams/system-context.svg)
- [Container Diagram](docs/diagrams/container.svg)
- [Component Diagram](docs/diagrams/component.svg)
- [Architecture Details](docs/architecture.md)
- [OpenAPI Specification](docs/openapi.yaml) (if available)

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15

### Development Tools

- [air](https://github.com/air-verse/air) - Live reload for Go applications
- [golangci-lint](https://golangci-lint.run/) - Go linters aggregator
- [gotestsum](https://github.com/gotestyourself/gotestsum) - Better test output formatter
- [sql-migrate](https://github.com/rubenv/sql-migrate) - Database migration tool

## Getting Started

1. Clone the repository

    ```bash
    git clone https://github.com/ngoctrng/bookz.git
    ```

2. Copy environment file and configure

    ```bash
    cp .env.example .env
    ```

3. Start dependencies

    ```bash
    make local-dev
    ```

4. Run database migrations

    ```bash
    make db/migrate
    ```

5. Start the HTTP server

    ```bash
    go run cmd/httpserver/main.go
    ```

6. Start the background worker

    ```bash
    go run cmd/worker/main.go
    ```

## Generating Swagger OpenAPI

This project uses [swaggo/swag](https://github.com/swaggo/swag) for API documentation.

1. Install swag if you haven't:

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. Generate Swagger docs:

    ```bash
    make swagger
    ```

3. The generated Swagger UI and OpenAPI spec will be available in the `docs/swagger` directory.

## Development

### Project Layout

- [cmd](http://_vscodecontentref_/0) - Entry points for executables (httpserver, worker, migrate)
- [internal](http://_vscodecontentref_/1) - Private application code
- [pkg](http://_vscodecontentref_/2) - Public shared packages
- [tools](http://_vscodecontentref_/3) - Development and deployment tools
- [docs](http://_vscodecontentref_/4) - Documentation and OpenAPI specs

### Testing

Run all tests:
```bash
make test
```

### Database Migrations

Create a new migration:

```bash
sql-migrate new -env="development" your-new-migration
```

### Development Tools

Hot reload during development:

```bash
make run
```

Run worker:

```bash
make worker
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License
