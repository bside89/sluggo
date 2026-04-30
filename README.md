# SlugGo

A URL shortener built in Go. You give it a long URL, it gives you a short one. That's it.

Behind the scenes: Gin handles the HTTP layer, PostgreSQL stores the URLs, Redis caches redirects so the database isn't hit on every click, and Snowflake IDs feed into HashIDs to generate the short codes.

## Stack

- **Go** — application
- **PostgreSQL** — persistent storage
- **Redis** — redirect cache (24h TTL by default)
- **Gin** — HTTP router
- **GORM** — database ORM
- **Swagger** — API docs, available at `/swagger`

## Prerequisites

- Go 1.25+
- Docker and Docker Compose

## Running locally

The local mode runs PostgreSQL and Redis in Docker, then starts the API process directly on your machine (useful for debugging).

**1. Copy the example env file:**

```bash
cp .env.example .env.local
```

Open `.env.local` and fill in at minimum:

```
DB_USER=your_user
DB_PASSWORD=your_password
HASH_SECRET_KEY=any_random_string
```

**2. Start the app:**

```bash
make local
```

This copies `.env.local` to `.env`, spins up the containers, waits for them to be healthy, and runs `go run ./cmd/api`.

The API will be at `http://localhost:8080`.

## Running with Docker

To build and run everything (app + database + cache) as containers:

```bash
make prod
```

## API

| Method | Path       | Description                  |
| ------ | ---------- | ---------------------------- |
| `POST` | `/shorten` | Create a short URL           |
| `GET`  | `/:hash`   | Redirect to the original URL |
| `GET`  | `/swagger` | API documentation            |

**Shorten a URL:**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/some/very/long/path"}'
```

Response:

```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

## Other commands

```bash
make build   # Compile binary to bin/sluggo
make docs    # Regenerate Swagger docs
```
