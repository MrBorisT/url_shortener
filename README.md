![CI](https://github.com/MrBorisT/url_shortener/actions/workflows/ci.yml/badge.svg)
# URL Shortener API

A small, complete backend API for creating and managing short links.

This project is intentionally scoped as a backend API slice: authentication, user-owned resources, PostgreSQL persistence, migrations, Docker setup, tests, CI, and a clear HTTP contract.

It is not trying to be a production-ready Bitly clone.

## Project goals

The goal of this project is to demonstrate backend engineering fundamentals in a compact service:

- JWT-based authentication
- user ownership and protected resources
- PostgreSQL data modeling
- database migrations
- validation and error handling
- public redirect behavior
- unit and integration testing
- Docker Compose local environment
- CI with tests and linting

## Features

- User registration
- Login with JWT
- Protected links API
- Users can manage only their own links
- Create short links
- List own links
- Get one own link
- Update original URL
- Disable link
- Delete link
- Public redirect by short code
- Click counter increments on redirect
- Disabled links return `410 Gone`
- Unknown short codes return `404 Not Found`
- URL validation
- PostgreSQL persistence
- SQL migrations
- Docker Compose setup
- Unit tests and integration tests
- GitHub Actions CI
- golangci-lint configuration

## Tech stack

- Go
- chi
- pgx
- PostgreSQL
- Docker / Docker Compose
- golang-migrate
- JWT
- bcrypt
- GitHub Actions
- golangci-lint

## API overview

| Method | Endpoint | Auth | Description |
|---|---|---:|---|
| `GET` | `/health` | No | Health check |
| `POST` | `/api/auth/register` | No | Register a new user |
| `POST` | `/api/auth/login` | No | Login and receive JWT |
| `POST` | `/api/links` | Yes | Create a short link |
| `GET` | `/api/links` | Yes | List current user's links |
| `GET` | `/api/links/{id}` | Yes | Get one current user's link |
| `PATCH` | `/api/links/{id}` | Yes | Update original URL |
| `DELETE` | `/api/links/{id}` | Yes | Delete link |
| `POST` | `/api/links/{id}/disable` | Yes | Disable link |
| `GET` | `/{short_code}` | No | Redirect to original URL |

Authenticated API endpoints return JSON errors:

```json
{
  "error": "link not found"
}
```

The public redirect endpoint may return plain HTTP errors, for example `404 Not Found` or `410 Gone`.

## Request examples

### Health check

```bash
curl http://localhost:8080/health
```

### Register

```bash
curl -i -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secret123"
  }'
```

### Login

```bash
curl -i -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secret123"
  }'
```

Response:

```json
{
  "token": "jwt_token_here"
}
```

For the examples below:

```bash
TOKEN="jwt_token_here"
```

### Create link

```bash
curl -i -X POST http://localhost:8080/api/links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "original_url": "https://example.com"
  }'
```

Example response:

```json
{
  "id": "0f38f54a-4a89-4d85-b6c7-9c763a777111",
  "user_id": "1e93f8e5-9a9c-4430-9873-dc43ad55c001",
  "original_url": "https://example.com",
  "short_code": "aB12xYz9",
  "click_count": 0,
  "created_at": "2026-04-30T18:00:00Z",
  "updated_at": "2026-04-30T18:00:00Z"
}
```

### List links

```bash
curl -i http://localhost:8080/api/links \
  -H "Authorization: Bearer $TOKEN"
```

### Get one link

```bash
curl -i http://localhost:8080/api/links/0f38f54a-4a89-4d85-b6c7-9c763a777111 \
  -H "Authorization: Bearer $TOKEN"
```

### Update original URL

```bash
curl -i -X PATCH http://localhost:8080/api/links/0f38f54a-4a89-4d85-b6c7-9c763a777111 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "original_url": "https://go.dev"
  }'
```

### Disable link

```bash
curl -i -X POST http://localhost:8080/api/links/0f38f54a-4a89-4d85-b6c7-9c763a777111/disable \
  -H "Authorization: Bearer $TOKEN"
```

After disabling, redirecting by its short code returns:

```http
410 Gone
```

### Delete link

```bash
curl -i -X DELETE http://localhost:8080/api/links/0f38f54a-4a89-4d85-b6c7-9c763a777111 \
  -H "Authorization: Bearer $TOKEN"
```

### Redirect

```bash
curl -i http://localhost:8080/aB12xYz9
```

Successful redirects return:

```http
302 Found
Location: https://example.com
```

Unknown short codes return:

```http
404 Not Found
```

Disabled links return:

```http
410 Gone
```

## Environment variables

The application reads configuration from environment variables.

| Variable | Description | Example |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | PostgreSQL database name | `urlshortener` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `postgres` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `APP_PORT` | HTTP server port | `8080` |
| `APP_ENV` | Runtime environment | `development` |
| `JWT_SECRET` | Secret used to sign JWTs | `change-me` |
| `JWT_TTL` | JWT lifetime | `24h` |

`JWT_SECRET` is required.

Docker Compose uses development credentials and a development JWT secret for local usage only. Do not reuse them in real deployments.

Example `.env` for local development:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=urlshortener
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
APP_PORT=8080
APP_ENV=development
JWT_SECRET=change-me-local-secret
JWT_TTL=24h
```

## Run with Docker Compose

Docker Compose starts:

- PostgreSQL
- migration container
- API container

```bash
docker compose up --build
```

The API will be available at:

```text
http://localhost:8080
```

Stop the stack:

```bash
docker compose down
```

Remove local database volume:

```bash
docker compose down -v
```

The Compose setup is intended for local development. It uses development database credentials and a local JWT secret.

## Run locally

Start PostgreSQL first. You can use Docker Compose for the database only:

```bash
docker compose up db
```

Create a local `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=urlshortener
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
APP_PORT=8080
APP_ENV=development
JWT_SECRET=change-me-local-secret
JWT_TTL=24h
```

Run migrations:

```bash
migrate -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable" \
  up
```

Run the API:

```bash
go run ./cmd/api
```

## Migrations

Migrations are stored in `migrations/`.

Current migrations:

```text
001_create_users_table
002_create_links_table
```

Run migrations up:

```bash
migrate -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable" \
  up
```

Rollback one migration:

```bash
migrate -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable" \
  down 1
```

## Testing

Run unit tests:

```bash
go test ./...
```

Run integration tests:

```bash
go test -tags=integration ./...
```

Integration tests require PostgreSQL and `TEST_DATABASE_URL`.

Example:

```bash
export TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/url_shortener_test?sslmode=disable"

migrate -path ./migrations \
  -database "$TEST_DATABASE_URL" \
  up

go test -tags=integration ./...
```

On Windows PowerShell:

```powershell
$env:TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/url_shortener_test?sslmode=disable"

migrate -path ./migrations `
  -database $env:TEST_DATABASE_URL `
  up

go test -tags=integration ./...
```

Test coverage includes:

- URL validation
- short code generation
- service behavior
- registration and login
- protected routes
- link ownership
- redirects
- disabled links
- click counter incrementing

## CI

GitHub Actions runs on pushes and pull requests to `main`.

The CI job:

1. Starts a PostgreSQL service
2. Installs `golang-migrate`
3. Runs migrations
4. Runs unit tests
5. Runs integration tests with the `integration` build tag
6. Runs `golangci-lint`

The linter configuration includes:

- `govet`
- `staticcheck`
- `ineffassign`
- `unused`
- `gofmt`

## Data model summary

### users

Stores registered users.

Main fields:

- `id` — UUID primary key
- `email` — unique email
- `password_hash` — bcrypt hash
- `created_at`
- `updated_at`

### links

Stores shortened links.

Main fields:

- `id` — internal UUID primary key
- `user_id` — owner ID, references `users(id)`
- `original_url` — destination URL
- `short_code` — public redirect code, unique
- `click_count` — incremented on successful redirect
- `disabled_at` — nullable timestamp
- `created_at`
- `updated_at`

The public short code is separate from the internal UUID. API management uses internal link IDs. Public redirects use short codes.

## Design decisions

### 302 redirect instead of 301

Redirects use `302 Found`, not `301 Moved Permanently`.

The target URL may change, and permanent redirect caching is undesirable for this MVP.

### Synchronous click counting

Click counting happens synchronously during redirect.

This keeps the system simple and visible for the current scope. A larger system could move this to an async event pipeline.

### Hard delete plus disable

The API supports both:

- hard delete for removing a link
- disable for keeping the link record while making the public redirect unavailable

Disabled redirects return `410 Gone`.

### Ownership enforced in storage queries

User-owned link operations include `user_id` in storage queries.

A user cannot read, update, disable, or delete another user's links. Non-owned links are treated as not found.

### Separate internal ID and public code

The internal link ID is a UUID.

The public short code is a separate unique value used only for redirects. This avoids exposing internal IDs as public redirect identifiers.

## Scaling notes

This project intentionally keeps the architecture simple. If it needed to handle higher traffic, the first areas to revisit would be:

- async click tracking
- caching hot redirects
- rate limiting
- short code collision strategy
- observability and structured logging
- database indexes based on real query patterns
- abuse prevention
- deployment and secret management

These are not implemented because the goal is a complete backend API slice, not a distributed URL shortening platform.

## Future improvements

Possible improvements:

- OpenAPI specification
- refresh tokens
- pagination for link listing
- custom aliases
- expiration dates
- rate limiting
- structured logging
- request IDs
- metrics
- better test coverage around edge cases
- deployment configuration
- frontend or simple admin UI

## Intentionally out of scope

The following are intentionally excluded from the current scope:

- Redis
- async workers
- distributed counters
- rate limiting
- frontend
- custom aliases
- OpenAPI
- analytics dashboard
- multi-region deployment
- production-grade secret management

The purpose of this project is to show a clean, finished backend service with realistic boundaries.
