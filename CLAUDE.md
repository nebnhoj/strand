# Strand

Go REST API built with Fiber, MongoDB, and JWT authentication.

## Commands

```bash
# Development (hot reload)
air

# Build
go build -o ./tmp/main .

# Run binary
./tmp/main

# Update Swagger docs
swag init

# Start local services (MongoDB, Redis)
docker-compose up
```

## Architecture

- **Framework:** Fiber v2
- **Database:** MongoDB (`configs/db.go`)
- **Auth:** JWT (HS512, 72h expiry) via `middlewares/jwt/`
- **Docs:** Swagger at `/api/docs`
- **Port:** 3000

### Module structure

Each module under `modules/` follows:
- `*.controller.go` — HTTP handlers, request parsing
- `*.model.go` — Structs, types
- `*.service.go` / `*.repository.go` — Business logic / DB queries

### Routes (`routes/routes.go`)

All routes prefixed with `/api`:

| Method | Path | Auth |
|--------|------|------|
| POST | /auth | public |
| POST | /users | JWT + admin |
| GET | /users | JWT + admin |
| GET | /users/:id | JWT + admin |
| POST | /todos | JWT |
| GET | /todos | JWT |

## Environment

Copy `.env.example` to `.env` and set:

```
MONGODB_URI=mongodb://localhost:27017
JWT_SECRET=your-secret
```

## Notes

- No test suite exists — test manually via Swagger or HTTP client
- Redis client is wired in but minimally used
- Pagination via `?page=1&limit=10&q=search` query params
