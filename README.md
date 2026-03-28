# Strand

A Go REST API boilerplate built with **Fiber v3**, **MongoDB**, **Redis**, and **JWT** authentication following **Clean Architecture** and **CQRS** design patterns.

## Stack

| Concern | Technology |
|---------|-----------|
| Framework | [Fiber v3](https://github.com/gofiber/fiber) |
| Database | MongoDB |
| Cache | Redis |
| Auth | JWT (HS512, 72h expiry) |
| Docs | Swagger UI (`/api/docs`) |
| Hot reload | [Air](https://github.com/cosmtrek/air) |

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Air (`go install github.com/air-verse/air@latest`)
- swag (`go install github.com/swaggo/swag/cmd/swag@latest`)

## Getting Started

```bash
# 1. Clone
git clone https://github.com/nebnhoj/strand.git
cd strand

# 2. Install dependencies
go mod tidy

# 3. Configure environment
cp .env.example .env

# 4. Start local services (MongoDB + Redis)
docker compose up mongodb redis -d

# 5. Run with hot reload
air
```

The API will be available at `http://localhost:3000`.
Swagger docs at `http://localhost:3000/api/docs`.

## Running with Docker

```bash
# Build and start all services (app + MongoDB + Redis)
docker compose up --build -d

# View logs
docker compose logs app -f
```

Ports when running via Docker Compose:

| Service | Host port |
|---------|-----------|
| API | 3001 |
| MongoDB | 27018 |
| Redis | 6380 |

## Seed Data

```bash
# Populate MongoDB with sample users and todos
go run ./cmd/seed/main.go
```

Default users:

| Email | Password | Role |
|-------|----------|------|
| admin@strand.dev | password | ADMIN |
| jane@strand.dev | password | USER |
| john@strand.dev | password | USER |

## Architecture

The project follows **Clean Architecture** with **CQRS**:

```
cmd/
  api/        — composition root (main.go)
  seed/       — database seeder
internal/
  domain/     — entities, repository interfaces, domain services
  application/
    auth/commands/    — authenticate command
    user/commands/    — create user command
    user/queries/     — list users, get user queries
    todo/commands/    — create todo command
    todo/queries/     — list todos query
  infrastructure/
    mongodb/   — MongoDB repository implementations
    jwt/       — JWT token service
    redis/     — Redis cache implementation
  interfaces/
    http/
      handlers/   — Fiber HTTP handlers
      middleware/  — JWT auth middleware
    router/     — route bindings
pkg/
  cache/      — Cache interface
  response/   — standard JSON response helpers
  validator/  — request validation
  apperrors/  — error message constants
docs/         — generated Swagger files
```

## API Endpoints

All routes are prefixed with `/api`.

| Method | Path | Auth | Permission | Description |
|--------|------|------|------------|-------------|
| POST | `/auth` | — | — | Authenticate, returns JWT |
| GET | `/users` | JWT | `users:read` | List users (paginated) |
| GET | `/users/:id` | JWT | `users:read` | Get user by ID |
| POST | `/users` | JWT | `users:write` | Create user |
| GET | `/todos` | JWT | `todos:read` | List todos (paginated) |
| POST | `/todos` | JWT | `todos:write` | Create todo |
| DELETE | `/cache` | JWT | `users:write` | Flush Redis cache |

### Pagination

Append query parameters to list endpoints:

```
?page=1&limit=10&q=search
```

## Auth Response

`POST /api/auth` returns:

```json
{
  "status": 200,
  "message": "success",
  "data": {
    "token": "eyJ...",
    "user_id": "uuid",
    "email": "admin@strand.dev",
    "permissions": ["users:read", "users:write", "todos:read", "todos:write"]
  }
}
```

## Permissions

Permissions are stored on each user document and embedded in the JWT at login.

| Role | Permissions |
|------|-------------|
| `ADMIN` | `users:read`, `users:write`, `todos:read`, `todos:write` |
| `USER` | `todos:read`, `todos:write` |

To add a role or change its permissions, edit `internal/domain/auth/permissions.go`.

## Caching

Read queries are cached in Redis (5 min TTL):

| Query | Key pattern |
|-------|-------------|
| List users | `users:list:{q}:{page}:{limit}` |
| Get user | `users:get:{id}` |
| List todos | `todos:list:{q}:{page}:{limit}` |

Cache is invalidated automatically on writes. To flush manually:

```bash
DELETE /api/cache
Authorization: Bearer <admin token>
```

## Development Commands

```bash
# Hot reload
air

# Build binary
go build -o ./tmp/main ./cmd/api

# Run binary
./tmp/main

# Regenerate Swagger docs
swag init -g cmd/api/main.go

# Seed database
go run ./cmd/seed/main.go
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```env
PORT=3000
APP_NAME=Strand API
APP_HEADER=Strand
JWT_SECRET=your-secret-here
MONGODB_URI=mongodb://strand:password@localhost:27017/strand?authSource=admin
REDIS_URL=redis://:password@localhost:6379
```
