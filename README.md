# Go Links (URL Shortener)

A simple internal "go links" service — short shortcuts (`go/oncall`, `go/design-system`)
that redirect to real URLs. Think of it as a mini self-hosted URL shortener for a team.

**Stack:** Go (net/http + chi) · SQLite · React + TypeScript + Vite + Tailwind

## Running locally

Requires Go 1.22+ and Node 18+.

```bash
# backend — http://localhost:8080
cd backend && go mod tidy && go run ./cmd/server

# frontend — http://localhost:5173
cd frontend && npm install && npm run dev


Backend tests: `cd backend && go test ./...`

## Assumptions

- Trusted internal users, no auth required — matches the real `go/` model
- Shortcuts are case-sensitive, no namespacing

## Tradeoffs

- **SQLite over Postgres** — zero infra to run, simpler for this scope
- **Create/delete only, no edit** — matched the requested API surface


## What I'd add next

- If time permits, I'd add the funtionality to sort the shortcuts either by the time of creation or alphabetical order and also would like to add some permissions as well.