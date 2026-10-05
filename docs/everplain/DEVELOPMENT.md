# Build and verify

Read upstream DEV_GUIDE.md. Use Go **1.27.0** from backend/go.mod and pnpm with the frozen frontend lockfile. Preserve the original architecture, all upstream capabilities and original standard/simple modes. Do not use production credentials or databases for local checks.

## Production asset build

```sh
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend typecheck
pnpm --dir frontend test:run --maxWorkers=1 --minWorkers=1 --no-file-parallelism
pnpm --dir frontend build
cd backend
go build -tags embed -trimpath -ldflags="-s -w" -o bin/everplain-gateway ./cmd/server
```

The frontend build outputs to backend/internal/web/dist. Use `pnpm --dir frontend exec vite build --config vite.config.ts --outDir dist` only for a standalone frontend artifact after type/i18n checks. This does not replace the full build checks.

For a memory-constrained build host, serialize processes and use bounded compiler parallelism (`GOMAXPROCS=1`, `go build -p 1`). These are build-time choices, not a feature-limited runtime mode or a production resource promise.

## Backend tests

Run upstream `make test-unit`, `make test-integration`, and golangci-lint under backend. The only current backend source difference is a setup DSN fix that reuses the existing database DSN builder to avoid treating an empty password as the next configuration keyword. It has focused regression coverage.

## Local runtime

Use upstream deploy/config.example.yaml and deployment documentation for a fresh isolated PostgreSQL database, Redis and administrator setup. Keep services bound to loopback for local development. No real provider/account secret is provided. Do not connect the disposable development database to production.

## Read-only UI preview

A separate downloadable preview bundle includes compiled Vue assets and tools/everplain/preview.py. It needs only Python 3, runs on loopback, marks all configuration as demonstrative, makes no Google/provider requests, returns unavailable usage, and rejects writes. It is not a runnable provider gateway and cannot complete OAuth. See tools/everplain/PREVIEW.md.

## Verification boundaries

Historical frontend checkpoint: 341 files / 2647 tests, typecheck, build and changed-file lint passed. Historical synthetic text-only load tests do not establish the resource use of this full-feature version. Cloud memory pressure has interrupted newer full build/type/suite attempts. Report final-head checks separately and do not merge on an invented pass.

A real Mac browser preview of the earlier bundle revealed dark-foreground, hardcoded table-header and overlay-banner issues. Source fixes use canonical tokens; a fresh build and browser review are required to accept the new full-feature bundle. No production deployment or real model call is authorized by these checks.
