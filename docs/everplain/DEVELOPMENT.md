# Build and verify the Everplain fork

## Base and resources

Read the upstream `DEV_GUIDE.md`. Its Go requirement is **1.27.0**, frontend package manager is **pnpm**, and the lockfile must remain frozen unless intentionally updated. Keep all development in a separate worktree and do not use a production database, provider account, or paid model API for tests.

The official Go 1.27.0 linux/amd64 archive was verified against the official go.dev release SHA256 `675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685`. This is a development toolchain, not an uploaded repository artifact.

## Build a single embedded server

```sh
pnpm --dir frontend install --frozen-lockfile
pnpm --dir frontend run build
cd backend
GOMAXPROCS=1 GOGC=20 GOMEMLIMIT=768MiB go build -p 1 -tags embed -trimpath -ldflags="-s -w" -o bin/everplain-gateway ./cmd/server
```

Vite writes its production assets to the existing `backend/internal/web/dist`; `-tags embed` places those assets inside the Go binary. No separate frontend runtime or new microservice is required. The upstream root `make build` compiles the backend before frontend and does not set `-tags embed`, so use the command order above or `make -f deploy/everplain/Makefile build` for this fork.

The lower compile parallelism and GC limit are for resource-constrained development hosts. The initial uncapped Ent package compile was killed by memory pressure; this does not establish the running gateway's RAM requirement. Do not confuse compiler memory with idle service memory.

## Tests

```sh
cd backend
go test ./internal/everplain ./internal/config
go test -tags unit ./internal/handler ./internal/server/... -run 'Everplain|SimpleMode|GroupModelAllowlist|ClientRequestID'
go test -tags unit ./internal/service -run 'SimpleModeRecordUsageWindowOptIn|OpenAIGatewayServiceRecordUsage_SimpleMode|ReserveInflightBalance_SimpleMode'
go test -tags unit ./...
go test -tags integration ./...
golangci-lint run ./...
```

Focused boundary tests use local fake handlers and miniredis. They do not contact model providers. Full upstream unit/integration/lint checks have their own dependency requirements; see the checked verification report for which stages actually ran. Never treat a focused pass as the full suite.

## Isolated local runtime

Use `deploy/everplain/config.example.yaml` as the lightweight starting profile. A fresh test instance needs a private `DATA_DIR`, its own PostgreSQL database, and its own Redis instance, all bound to loopback. Supply local-only database/Redis parameters through environment variables. `CONFIG_FILE` selects the example or a private copy. `AUTO_SETUP=true` plus local-only administrator settings can initialize an empty isolated database. Do not put actual secrets into committed YAML or shell examples.

Production deployment is deliberately out of scope. Before deployment, use TLS, authorized provider credentials, tested model mappings, a dedicated Everplain service key, backups, log/usage retention, and separately validated performance/security/recovery tests. The source tree preserves upstream distribution and release workflows; do not trigger a tag release without reviewing fork targets, names, signing, licensing and deployment intent.

## Synthetic provider fixture

`tools/everplain/seed-local-test.sql` refuses to run outside the explicitly disposable `everplain_gateway_test` database. It creates only a synthetic service user/key and one localhost Anthropic adapter with an `everplain-test` model mapping. Supply the psql variables `service_key` and `upstream_key` from private, newly generated runtime values. No authentication value is committed. Set `EVERPLAIN_TEST_SERVICE_KEY` only in the isolated test process when invoking the smoke script.

The mock supports `--tls-cert <temporary-cert> --tls-key <temporary-key>` and binds only `127.0.0.1:18082`. Use a short-lived test certificate with an IP SAN for `127.0.0.1`, set `SSL_CERT_FILE` for the isolated gateway process only, and give that process a private config with the URL allowlist enabled, upstream hosts limited to `127.0.0.1`, and private hosts permitted. Do not change the system CA store or disable certificate verification. When the allowlist is enabled upstream enforces HTTPS regardless of `allow_insecure_http`.

Initialize the isolated gateway at `127.0.0.1:18081`, apply the fixture, and restart it so the existing scheduler sees the account. Run `python tools/everplain/smoke-local.py` while upstream is disabled. Then explicitly enable the synthetic adapter in that isolated process and run `python tools/everplain/smoke-local.py --configured`. These scripts never contact an external model service. The mock injects a fixed 25 ms delay and emits 10 input/5 output tokens; synthetic latency is not a model-provider SLO.

Stop the gateway, mock, PostgreSQL and Redis processes after validation. Do not reuse test keys, certificates or databases in a deployment.
