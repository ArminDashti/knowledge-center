# Smoke test knowledge center mvp

| Field | Value |
|-------|-------|
| Agent | Auto |
| Date | 2026-08-07 |
| Device | ARMIN-DESKTOP |

## Transcript

### User
Approved (Go + npm installs)

### Agent
- Postgres healthy on 5434
- API runs on :8091 (8080/8090 occupied by Docker/other)
- WebUI Vite on :5173 proxying /api -> 8091
- Smoke: health ok; import https://example.com -> success + page title Example Domain; /list /import /scraping return 200
- cursor-ide-browser MCP unavailable; API/proxy verified; Windows DOM scrape flaky
