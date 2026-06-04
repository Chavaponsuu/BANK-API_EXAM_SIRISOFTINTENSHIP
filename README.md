# Go Gin CRUD API

RESTful CRUD API built with Go and Gin framework, connected to PostgreSQL.

## Tech Stack

- **Go 1.23+** with **Gin** HTTP framework
- **PostgreSQL 16** via Docker Compose
- **golang-migrate** for database migrations
- **swaggo/gin-swagger** for API documentation

## Project Structure

```
├── main.go                 # Entry point
├── Makefile                # Build & dev commands
├── docker-compose.yml      # PostgreSQL container
├── config/config.go        # Environment config loader
├── db/db.go                # Database connection & migration runner
├── models/example.go       # Database models & request structs
├── repositories/example.go # SQL query layer
├── handlers/example.go     # HTTP handlers
├── dto/
│   ├── response.go         # Base response wrapper (standardized JSON)
│   └── example.go          # Example DTOs & model converters
├── routes/routes.go        # Route definitions
├── middleware/middleware.go # Logger, CORS, panic recovery
├── migrations/             # Versioned SQL migration files
└── docs/                   # Generated Swagger docs
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose

### Quick Start

```bash
# Clone the repository
git clone https://github.com/krizad/go-gin-api.git
cd go-gin-api

# Start PostgreSQL and the API
make dev
```

The API will be available at `http://localhost:8080`.

### Step by Step

```bash
# 1. Start PostgreSQL
make db-up

# 2. Run database migrations
make migrate-up

# 3. Start the API
make run
```

## API Endpoints

| Method   | Endpoint                            | Description         |
| -------- | ----------------------------------- | ------------------- |
| `GET`    | `/api/v1/examples`                  | List all examples   |
| `POST`   | `/api/v1/examples`                  | Create an example   |
| `GET`    | `/api/v1/examples/search?email=x`   | Search by email     |
| `GET`    | `/api/v1/examples/:id`              | Get by ID           |
| `PUT`    | `/api/v1/examples/:id`              | Full update         |
| `PATCH`  | `/api/v1/examples/:id`              | Partial update      |
| `DELETE` | `/api/v1/examples/:id`              | Delete              |
| `GET`    | `/health`                           | Health check        |

### Swagger UI

Visit `http://localhost:8080/swagger/index.html` for interactive API documentation.

## Standard API Response

All endpoints return a standardized response format:

```json
{
  "success": true,
  "message": "OK",
  "data": { ... },
  "errors": null,
  "meta": { "page": 1, "per_page": 20, "total": 100 }
}
```

### Example Requests

**Create an example:**

```bash
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

**List all examples:**

```bash
curl http://localhost:8080/api/v1/examples
```

**Get by ID:**

```bash
curl http://localhost:8080/api/v1/examples/1
```

**Partial update:**

```bash
curl -X PATCH http://localhost:8080/api/v1/examples/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'
```

**Delete:**

```bash
curl -X DELETE http://localhost:8080/api/v1/examples/1
```

## Makefile Commands

| Command              | Description                        |
| -------------------- | ---------------------------------- |
| `make run`           | Start the API server               |
| `make build`         | Build the binary to `bin/`         |
| `make test`          | Run tests with race detection      |
| `make lint`          | Run `go vet`                       |
| `make swagger`       | Regenerate Swagger docs            |
| `make db-up`         | Start PostgreSQL container         |
| `make db-down`       | Stop PostgreSQL container          |
| `make db-reset`      | Reset database (drop + recreate)   |
| `make migrate-up`    | Apply all pending migrations       |
| `make migrate-down`  | Rollback all migrations            |
| `make dev`           | Full dev setup (db + migrate + api)|
| `make clean`         | Remove build artifacts             |

## Environment Variables

Copy `.env` and adjust as needed:

| Variable         | Default      | Description                |
| ---------------- | ------------ | -------------------------- |
| `DB_HOST`        | `localhost`  | PostgreSQL host            |
| `DB_PORT`        | `5432`       | PostgreSQL port            |
| `DB_USER`        | `postgres`   | Database user              |
| `DB_PASSWORD`    | `postgres`   | Database password          |
| `DB_NAME`        | `go_gin_api` | Database name              |
| `DB_SSLMODE`     | `disable`    | SSL mode                   |
| `APP_PORT`       | `8080`       | API listen port            |
| `AUTO_MIGRATE`   | `true`       | Auto-run migrations on start|

## Database Migrations

Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate) and embedded in the binary.

```bash
# Run pending migrations
make migrate-up

# Rollback all migrations
make migrate-down

# Create a new migration file
# Add 000002_*.up.sql and 000002_*.down.sql to migrations/
```

## License

MIT
