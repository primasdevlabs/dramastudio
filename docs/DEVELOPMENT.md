# Development Guide

Setup, run, and test instructions for DramaStudio — from a fresh clone to a
running stack with verified tests.

## Prerequisites

| Tool | Version | Required for |
|---|---|---|
| Go | 1.22+ | backend (`apps/api`, `apps/worker`, `apps/scheduler`) |
| Node.js | 20+ | frontend (`apps/web`) |
| FFmpeg | any recent | actual video rendering only — API and tests run without it |
| PostgreSQL | 14+ | only when `DRAMASTUDIO_STORE=postgres` |
| Temporal | — | only when `WORKFLOW_ENGINE=temporal` |
| S3-compatible store | — | only when `STORAGE_BACKEND=s3` |

A default development run needs **only Go + Node.js**: SQLite file storage,
in-process workflow engine, and local object storage are the defaults.

## 1. Install

```bash
# from the repo root
go mod download          # backend dependencies
cd apps/web
npm ci                   # frontend dependencies (package-lock.json)
cd ../..
```

## 2. Configure environment

Two env files ship as templates:

- `.env.example` → `.env` (repo root — backend config)
- `apps/web/.env.local.example` → `apps/web/.env.local` (frontend config)

Both `.env` files are already created for development with working defaults.
The committed `.env` in this repo sets:

```dotenv
DRAMASTUDIO_ENV=development   # sqlite + local engine + local storage
PORT=9471                     # API port
REQUIRE_AUTH=false            # injects a dev principal — no login needed
```

### Loading `.env` into the Go processes

The Go binaries read **real environment variables** — there is no `.env`
autoload. Load it before starting a process:

**Git Bash / bash:**

```bash
set -a; source .env; set +a
go run ./apps/api
```

**PowerShell:**

```powershell
Get-Content .env | Where-Object { $_ -match '^\s*[^#].*=' } | ForEach-Object {
  $k,$v = $_ -split '=',2; [Environment]::SetEnvironmentVariable($k.Trim(),$v.Trim(),'Process')
}
go run .\apps\api
```

Because every value in the committed `.env` matches a built-in default, you can
also just run `go run ./apps/api` with **no env at all** — it works out of the
box. Set variables only for the parts you want to change; see
`.env.example` for the full annotated list.

The frontend reads `.env.local` automatically (`NEXT_PUBLIC_API_URL`).

## 3. Run

### API server — terminal 1

```bash
go run ./apps/api
```

Starts on `http://localhost:9471`. On first boot it creates
`data/dramastudio.db`, applies `migrations/`, and prepares `data/objects` for
local media storage.

Health check:

```bash
curl http://localhost:9471/health
```

### Frontend — terminal 2

```bash
cd apps/web
npm run dev
```

Starts on `http://localhost:8742` and proxies API calls to
`NEXT_PUBLIC_API_URL` (`http://localhost:9471` by default).

With `REQUIRE_AUTH=false` every request runs as the dev principal — skip login
and open <http://localhost:8742/projects> directly. Set `REQUIRE_AUTH=true` +
`JWT_SECRET` to exercise the real `/v1/auth/login` flow through `/login`.

### Seed a known login (auth-enabled mode)

```bash
go run ./apps/seed
```

Creates the dev organization and an owner account — idempotent, safe to
re-run. Defaults (override with `SEED_ORG`, `SEED_NAME`, `SEED_EMAIL`,
`SEED_PASSWORD`):

```
email:    dev@dramastudio.local
password: dramastudio-dev
```

Sign in at `/login`, or create a new organization at `/register`.

### Optional processes

```bash
go run ./apps/scheduler   # polls publishing for due publications (default 30s)
go run ./apps/worker      # temporal activity workers — exits immediately on local engine
```

You only need `scheduler` to watch scheduled publications go out, and `worker`
only when `WORKFLOW_ENGINE=temporal`.

## 4. Test

### Backend

```bash
go build ./...    # compile every package
go vet ./...      # static analysis
go test ./...     # full suite — includes module-level tests in ./tests
```

Module tests live in `tests/` (`identity_test.go`, `media_test.go`,
`publishing_test.go`, …) and run against the in-memory/SQLite adapters — no
external services needed. Mock AI providers make generation paths deterministic.

### Frontend

```bash
cd apps/web
./node_modules/.bin/tsc --noEmit   # typecheck
npm run build                      # production build (typecheck + bundle)
```

### Smoke test — end-to-end

With the API running, verify the whole pipeline surface:

```bash
# health
curl http://localhost:9471/health

# create a project
curl -X POST http://localhost:9471/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"name":"Smoke Test","genre":"drama","language":"en","mode":"monitored"}'

# list projects (returns {"items":[...]})
curl http://localhost:9471/v1/projects
```

With auth enabled (`REQUIRE_AUTH=true`), register and log in first:

```bash
curl -X POST http://localhost:9471/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"org_name":"Dev Org","email":"dev@example.com","password":"password123","name":"Dev"}'

curl -X POST http://localhost:9471/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@example.com","password":"password123"}'
# → {"access_token":"...","token_type":"Bearer","user":{...}}
# pass the token as  -H "Authorization: Bearer <token>"
```

## 5. Switching to production-shaped services

| Component | Dev default | Production |
|---|---|---|
| Database | `DRAMASTUDIO_STORE=sqlite` | `postgres` + `DATABASE_URL` |
| Workflow | `WORKFLOW_ENGINE=local` | `temporal` + `TEMPORAL_HOST_PORT` + `apps/worker` |
| Storage | `STORAGE_BACKEND=local` | `s3` + endpoint/bucket/keys |
| Auth | `REQUIRE_AUTH=false` | `true` + `JWT_SECRET` |
| Tracing | `OTEL_TRACE_EXPORTER=none` | `stdout` (or wire an OTLP exporter) |

`configs/config.go` fails fast on missing required variables for the selected
profile — misconfiguration surfaces at boot, not at runtime.
