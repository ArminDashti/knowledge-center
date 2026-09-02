# Knowledge Center API

Gin + PostgreSQL backend for importing websites, scraping with per-domain CSS profiles, and storing page content.

## Requirements

- Go 1.22+
- PostgreSQL 16 (local Docker compose provided)

## Configuration

Copy `.env.example` to `.env`:

```
DATABASE_URL=postgres://kc:kc@127.0.0.1:5434/knowledge_center?sslmode=disable
HTTP_ADDR=:8091
```

## Run

```powershell
# start Postgres (after you approve Docker pull/run)
.\.armin\docker-scripts\run-on-docker-local.ps1

go run ./cmd/server
```

## Endpoints

- `GET /api/health`
- `POST /api/imports` — `{ "url": "https://example.com" }`
- `GET /api/imports`
- `GET /api/imports/:id`
- `GET /api/imports/:id/pages/:pageId`
- `GET|POST /api/scrape-profiles`
- `GET|PUT|DELETE /api/scrape-profiles/:id`
