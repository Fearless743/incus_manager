# AGENTS.md

Incus Manager: Go API + React/Vite UI for multi-user Incus host and instance management (PostgreSQL, JWT, per-host Incus REST clients).

## Layout

| Path | Role |
|------|------|
| `backend/cmd/server.go` | Entry point: DB init, GORM AutoMigrate, route registration, static file serving |
| `backend/internal/handler/` | HTTP handlers |
| `backend/internal/service/` | Business logic; `incus_api.go` wraps Incus REST API |
| `backend/internal/model/models.go` | GORM models (`User`, `Host`, `Instance`) |
| `frontend/src/` | React 18 JSX (not TypeScript): pages, components, `services/api.js` |
| `API_DOCUMENTATION.txt` | API reference (prefer over guessing routes) |

Go module import path: `incus-manager/...` (hyphenated). Routing uses stdlib `http.ServeMux`, not a web framework.

## Commands

```bash
# Full production build → ./incus-manager at repo root
./build.sh

# Local dev (two terminals; needs PostgreSQL)
cp backend/.env.example backend/.env   # then export vars or set in shell — Go does NOT load .env files
cd backend && go mod tidy && go run cmd/server.go   # :8080
cd frontend && npm install && npm run dev           # :5173

# Docker (postgres + app)
docker-compose up --build -d   # app on :8080

# Frontend lint only (no test script, no Go tests in repo)
cd frontend && npm run lint
```

CI (`.github/workflows/docker-publish.yml`) only builds/pushes `Dockerfile.multi` on main — no lint or test gates.

## Dev vs production serving

- **Dev:** run frontend and backend separately. Vite proxies `/api` and `/ws` to `:8080` (`frontend/vite.config.js`), but `frontend/.env.development` sets `VITE_API_URL=http://localhost:8080/api`, so axios calls the backend directly.
- **Production/Docker:** single binary serves API + SPA from **`/root/dist`** (hardcoded in `server.go`).
- **`./build.sh` copies `frontend/dist` → `backend/dist`**, but the binary still reads `/root/dist`. The local `./incus-manager` binary is API-only unless you run in Docker or symlink/copy dist to `/root/dist`. Do not assume `./build.sh && ./incus-manager` serves the UI locally.

## Environment

Backend reads **environment variables only** (`backend/internal/config/config.go`):

- `DATABASE_URL` — required (default in code omits host port; use `localhost:5432` locally)
- `JWT_SECRET` — required in production
- `PORT` — default `8080`

`INCUS_URL` / `INCUS_CERT` in `.env.example` are **unused**. Incus connectivity is per-host: address + certificate PEM stored in the `Host` row and passed to `IncusServiceFactory.GetClient`.

## Database

- GORM `AutoMigrate` runs on every server start in `server.go` (`User`, `Host`, `Instance`).
- No migration files or separate migrate command.
- `backend/pkg/database/` exists but is **not imported** anywhere — dead code; do not wire new code through it unless intentionally reviving it.

## Incus integration (non-obvious)

- Each host gets an Incus **project** name (`host-{name}-{userID}`) for isolation.
- `NewIncusClient(address, certificate, "")` passes the stored PEM as cert with an empty key; mTLS only works if the PEM is a valid cert+key pair parseable by `tls.X509KeyPair`. Otherwise TLS uses `InsecureSkipVerify: true`.
- Instance sharing is stored as JSON in `Instance.SharedWith`, not a separate table.
- WebSocket hub is mounted at `/ws` but the frontend does not consume it yet.

## Conventions

- Frontend: plain `.jsx`, ESLint 10 flat config (`frontend/eslint.config.js`).
- Handler auth: JWT via `middleware.Authenticate`; user ID from request context.
- When changing API routes, update `backend/cmd/server.go`, handlers, `frontend/src/services/api.js`, and `API_DOCUMENTATION.txt` together.
- Docker Go build uses `GOPROXY=https://goproxy.cn,direct` in `Dockerfile.multi` (China mirror); local `go build` does not.
