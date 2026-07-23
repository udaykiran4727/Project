# Go Links (URL Shortner)

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

## API

| Method | Path | Description |
|---|---|---|
| POST | `/api/links` | Create a link |
| GET | `/api/links` | List all links |
| GET | `/api/links/:id` | Get one link |
| DELETE | `/api/links/:id` | Delete a link |
| GET | `/:shortcut` | Redirects to the destination |


## What I'd add next

- Auth & permissions
- Edit support
- Pagination on link list
- CI (Go tests + TypeScript build) on push
- Frontend tests