# Employee Management System

A modular, scalable backend system for managing employee data, attendance, payroll, overtime, reimbursements, and audit logs. Built with Go, following clean architecture principles.

## Features

- User authentication and authorization
- Attendance tracking
- Overtime management
- Payroll processing
- Reimbursement requests
- Audit logging
- RESTful API with Swagger documentation
- Redis caching for improved performance

## Project Structure

```
├── cmd/                # Application entrypoints (main, API, seed)
├── config/             # Configuration files (YAML)
├── database/           # Migrations and seed data
├── docs/               # API documentation (Swagger)
├── internal/           # Core business logic (attendance, payroll, etc.)
│   ├── <module>/       # Each domain module (repository, usecase, entity, etc.)
├── pkg/                # Shared utilities and helpers
├── main.go             # Main application entrypoint
├── Makefile            # Common development commands
├── go.mod, go.sum      # Go modules
```

## Architecture

This project follows the principles of **Clean Architecture** to ensure separation of concerns, testability, and scalability. The main layers are:

- **Entities**: Core business models and logic, independent of frameworks and external systems (`internal/<module>/entity/`).
- **Use Cases**: Application-specific business rules, orchestrating entities and repositories (`internal/<module>/usecase/`).
- **Repositories (Interfaces & Implementations)**: Abstractions and concrete implementations for data access (`internal/<module>/repository/`).
- **Delivery (Controllers/Handlers)**: Interfaces for interacting with the outside world, such as HTTP handlers (`internal/<module>/delivery/`).
- **Infrastructure**: Shared utilities, database helpers, and middleware (`pkg/`, `config/`, `database/`).

This structure allows each module (attendance, payroll, etc.) to be developed, tested, and maintained independently, promoting a modular and robust backend system.

## Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL (or your configured DB)
- Redis

### Setup

1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd employee-management
   ```
2. Copy and edit configuration:
   ```sh
   cp config/api/config.example.yml config/api/config-local.yml
   # Edit config-local.yml as needed
   ```
3. Run database migrations:
   ```sh
   make migrate-up
   # You can override DB connection variables if needed, e.g.:
   # make migrate-up DB_USER=admin DB_PASS=secret DB_NAME=employee_management
   ```
4. Seed initial data (optional):
   ```sh
   go run main.go seed
   ```
5. Start the API server:
   ```sh
   go run main.go http
   ```

### Running Unit Tests

To run all unit tests:

```sh
go test ./...
```

You can also run tests for a specific package, for example:

```sh
go test ./internal/attendance/...
```

## Makefile Commands

The Makefile provides helpful commands for development and database management. You can override database connection variables as needed (e.g., `make migrate-up DB_USER=myuser DB_PASS=mypass`).

- `make migrate-up` — Run all database migrations (up).
- `make migrate-down` — Roll back the latest database migration (down).
- `make migrate-create name=your_migration_name` — Create a new migration file with the given name.
- `make mock-repo domain=yourdomain` — Generate repository mocks for a specific domain.
- `make mock-usecase domain=yourdomain` — Generate usecase mocks for a specific domain.

You can also set database connection variables:

- `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_SSLMODE`

Example:

```sh
make migrate-up DB_USER=admin DB_PASS=secret DB_NAME=employee_management
```

## API Documentation

- Swagger UI available at `/swagger` endpoint when the server is running.
- See `docs/swagger.yaml` for the OpenAPI spec.

## Performance Testing

This project includes a k6 load testing script (`k6.js`) to benchmark API performance and concurrency.

Example k6 report:

![k6 Report](https://imgur.com/a/f8MKiUS)

Resource usage visualization (memory & CPU):

![Resource Usage](https://imgur.com/a/XLPRyxZ)

To run the test:

```sh
k6 run k6.js
```

You can adjust the number of virtual users and duration in the script as needed.

## TODO

- [ ] Implement cursor-based pagination.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes
4. Push to the branch (`git push origin feature/YourFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
