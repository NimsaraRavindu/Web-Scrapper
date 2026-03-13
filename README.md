# Web Scrapper API

A Go-based RSS aggregation backend that lets users register, add RSS feeds, follow feeds, and read the latest posts from feeds they follow.

The service exposes a REST API, stores data in PostgreSQL, and runs a background scraper that continuously fetches and ingests RSS items.

## Features

- User creation with generated API keys
- API key authentication with `Authorization: ApiKey <key>`
- Create and list RSS feeds
- Follow and unfollow feeds
- Fetch latest posts for the authenticated user
- Background RSS scraper with configurable concurrency and interval (currently hardcoded in `main.go`)
- SQL-first data layer using sqlc-generated queries

## Tech Stack

- Go
- Chi router + CORS middleware
- PostgreSQL
- sqlc (query code generation)
- Goose-style SQL migrations

## Project Structure

- `main.go`: server bootstrap, routes, DB connection, scraper startup
- `handler_*.go`: HTTP handlers
- `middleware_auth.go`: API key auth middleware
- `scrapper.go`: background feed scraping loop
- `rss.go`: RSS fetch + XML parsing
- `sql/schema`: database migrations
- `sql/queries`: SQL queries used by sqlc
- `internal/database`: generated query code

## Prerequisites

- Go (version compatible with `go.mod`)
- PostgreSQL
- Optional tools:
  - goose (for migrations)
  - sqlc (to regenerate `internal/database` if SQL changes)

## Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
DB_URL=postgres://postgres:postgres@localhost:5432/webscrapper?sslmode=disable
```

## Database Setup

1. Create a PostgreSQL database.
2. Apply SQL migration files from `sql/schema` in order:
   - `001_users.sql`
   - `002_users_apikey.sql`
   - `003_feeds.sql`
   - `004_feed_follows.sql`
   - `005_feeds_last_fetched_at.sql`
   - `006_posts.sql`

If you use goose, point it at `sql/schema` and your database URL.

## Run the API

```bash
go run .
```

Server starts on `PORT` and serves routes under `/v1`.

## Background Scraper

A scraper goroutine starts automatically on boot:

- Concurrency: `10`
- Interval: every `1 minute`

It selects feeds by oldest `last_fetched_at`, fetches RSS items, and inserts new posts (duplicate post URLs are ignored).

## API Endpoints

Base path: `/v1`

### Public

- `GET /healthz` - readiness check
- `GET /err` - test error response
- `POST /users` - create user
- `GET /feeds` - list feeds

### Authenticated

Requires header:

```http
Authorization: ApiKey <your_api_key>
```

- `GET /users` - get current user
- `POST /feeds` - create feed (also stores owner)
- `POST /feedfollows` - follow a feed
- `GET /feedfollows` - list feeds current user follows
- `DELETE /feedfollows/{feedFollowID}` - unfollow
- `GET /posts` - get latest posts for current user (limit 10)

## Example Requests

### 1) Create user

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"nimsara"}'
```

### 2) Create feed (authenticated)

```bash
curl -X POST http://localhost:8080/v1/feeds \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <API_KEY>" \
  -d '{"name":"Hacker News","url":"https://hnrss.org/frontpage"}'
```

### 3) Follow a feed

```bash
curl -X POST http://localhost:8080/v1/feedfollows \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey <API_KEY>" \
  -d '{"feed_id":"<FEED_UUID>"}'
```

### 4) Get posts

```bash
curl -X GET http://localhost:8080/v1/posts \
  -H "Authorization: ApiKey <API_KEY>"
```

## Notes

- Feed URLs are unique.
- A user cannot follow the same feed twice.
- Posts are unique by URL.
- On startup, `main.go` currently performs a sample fetch from `https://wagslane.dev/index.xml` before server boot.

## Development

### Regenerate sqlc code

```bash
sqlc generate
```

### Run tests

```bash
go test ./...
```

(If no tests are present yet, this will simply compile-check packages.)
