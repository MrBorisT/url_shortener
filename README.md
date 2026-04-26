# URL Shortener API

A small backend URL shortener written in Go.

It supports user registration, JWT login, user-owned short links, public redirects, link disabling, click counting, PostgreSQL persistence, migrations, and Docker Compose setup.

## What it does

This API allows users to create and manage their own shortened links.

Authenticated users can:

- create short links
- list their links
- get a single link
- update the original URL
- disable a link
- delete a link

Public users can open a short URL and get redirected to the original URL.

Example:

```text
GET /abc123XY
```

redirects to the original URL stored for that short code.

## Features

- JWT auth
- user-owned links
- redirect by short code
- disable links
- click counter
- PostgreSQL persistence
- database migrations
- Docker Compose setup

## Tech stack

- Go
- chi
- pgx
- PostgreSQL
- Docker
- Docker Compose
- golang-migrate
- JWT
- bcrypt

## API

Base URL:

```text
http://localhost:8080
```

### Auth

| Endpoint | Method | Auth | Description |
|---|---:|---:|---|
| `/api/auth/register` | POST | No | Register a new user |
| `/api/auth/login` | POST | No | Login and receive a JWT token |

### Links

| Endpoint | Method | Auth | Description |
|---|---:|---:|---|
| `/api/links/` | POST | Yes | Create a short link |
| `/api/links/` | GET | Yes | List current user's links |
| `/api/links/{id}` | GET | Yes | Get one link by ID |
| `/api/links/{id}` | PATCH | Yes | Update the original URL |
| `/api/links/{id}` | DELETE | Yes | Delete a link |
| `/api/links/{id}/disable` | POST | Yes | Disable a link |

### Public

| Endpoint | Method | Auth | Description |
|---|---:|---:|---|
| `/health` | GET | No | Health check |
| `/{short_code}` | GET | No | Redirect to original URL |

## Request examples

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

### Create link

```bash
curl -i -X POST http://localhost:8080/api/links/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer jwt_token_here" \
  -d '{
    "original_url": "https://example.com"
  }'
```

### List links

```bash
curl -i http://localhost:8080/api/links/ \
  -H "Authorization: Bearer jwt_token_here"
```

### Get one link

```bash
curl -i http://localhost:8080/api/links/{id} \
  -H "Authorization: Bearer jwt_token_here"
```

### Update link

```bash
curl -i -X PATCH http://localhost:8080/api/links/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer jwt_token_here" \
  -d '{
    "original_url": "https://example.org"
  }'
```

### Disable link

```bash
curl -i -X POST http://localhost:8080/api/links/{id}/disable \
  -H "Authorization: Bearer jwt_token_here"
```

### Delete link

```bash
curl -i -X DELETE http://localhost:8080/api/links/{id} \
  -H "Authorization: Bearer jwt_token_here"
```

### Redirect

```bash
curl -i http://localhost:8080/{short_code}
```

A valid active short code returns a `302 Found` redirect.

## Run with Docker

Start the full stack:

```bash
docker compose up --build
```

This starts:

- PostgreSQL
- migration container
- API container

The API will be available at:

```text
http://localhost:8080
```

Stop containers:

```bash
docker compose down
```

Stop containers and remove the database volume:

```bash
docker compose down -v
```

## Environment variables

The app reads configuration from environment variables.

| Variable | Required | Default | Description |
|---|---:|---:|---|
| `DB_HOST` | Yes | - | PostgreSQL host |
| `DB_PORT` | Yes | - | PostgreSQL port |
| `DB_NAME` | Yes | - | PostgreSQL database name |
| `DB_USER` | Yes | - | PostgreSQL user |
| `DB_PASSWORD` | Yes | - | PostgreSQL password |
| `DB_SSLMODE` | Yes | - | PostgreSQL SSL mode |
| `APP_PORT` | No | `8080` | HTTP server port |
| `APP_ENV` | No | - | Use `production` in Docker to skip loading `.env` |
| `JWT_SECRET` | Yes | - | Secret key used to sign JWT tokens |
| `JWT_TTL` | No | `24h` | JWT lifetime |

Docker Compose already provides the required variables for local container usage.

For local non-Docker runs, create a `.env` file or export the variables manually.

Example:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=urlshortener
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable
APP_PORT=8080
JWT_SECRET=change_me
JWT_TTL=24h
```

## Migrations

Migrations are stored in the `migrations/` directory.

Current migrations create:

- `users` table
- `links` table
- index on `links.user_id`

When using Docker Compose, migrations run automatically through the `migrate` service.

Manual migration example:

```bash
migrate \
  -path migrations \
  -database "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable" \
  up
```

Rollback example:

```bash
migrate \
  -path migrations \
  -database "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable" \
  down
```

## Testing

Run Go tests:

```bash
go test ./...
```

For a basic manual smoke test:

```bash
docker compose up --build
```

Then check:

```bash
curl -i http://localhost:8080/health
```

Recommended manual flow:

1. Register a user
2. Login and copy the JWT token
3. Create a link
4. Open the returned short code
5. Check that `click_count` increases
6. Disable the link
7. Confirm that redirect returns `410 Gone`
8. Delete the link

## Design decisions

### 302 redirect

Redirects use `302 Found`.

This is intentional for an MVP because the target URL can be changed later. A permanent `301` redirect would be more aggressive and may be cached by clients.

### Sync click counting

Click counting is done synchronously during redirect.

This keeps the implementation simple and transparent. For a small portfolio backend this is acceptable. In a higher-load system, this could be moved to a queue, background worker, Redis, or batched writes.

### Hard delete + disable

The API supports both:

- hard delete through `DELETE /api/links/{id}`
- soft disabling through `POST /api/links/{id}/disable`

Disabled links remain in the database but no longer redirect.

Deleted links are removed from storage.

### No Redis / rate limiting / frontend by design

This project intentionally keeps the scope small.

Not included by design:

- Redis cache
- async workers
- rate limiting
- frontend UI
- analytics dashboard
- custom aliases

The goal is a focused backend API, not a full Bitly clone.

## Future improvements

- Add automated tests for handlers, services, validation, and storage
- Add refresh tokens
- Add rate limiting for auth and redirects
- Add custom short codes
- Add pagination for link listing
- Add structured logging
- Add OpenAPI / Swagger documentation
- Add request ID middleware
- Add graceful shutdown
- Add CI pipeline
- Add Makefile for common commands
- Add soft delete with `deleted_at` instead of hard delete
- Add redirect analytics by date, IP hash, or user agent
