# Unified management verification

2026-10-05. Source candidate, not production acceptance.

## Standard CI evidence

The unchanged source head `3f386af1907884ac0421653afec93c95d4c7cabf` completed [CI 37289065973](https://github.com/huyanxius/everplain-gateway/actions/runs/37289065973) and [Security Scan 37289066158](https://github.com/huyanxius/everplain-gateway/actions/runs/37289066158) successfully. This includes full Go unit/integration suites, golangci-lint, complete frontend tests/build, shell checks, release helpers, govulncheck and frontend dependency audit. The previously recorded cloud-only DNS/Unix-socket failures below do not reproduce on that standard runner. No original DNS fixture, SSRF guard or plugin test needs removal.

The release workflow and Dockerfile already contain embedded-frontend production build paths. The candidate's added read-only CI step reuses the existing frontend output and Go cache to compile a Linux/amd64 `-tags embed` gateway and save its checksum/build artifact, without publishing a release/image or acquiring write permissions. It does not repeat frontend build or tests. That additional exact binary build has not run until the new commit's CI completes; the original green run does not cover a later workflow edit.

## Preservation and source recovery

The restored `beabb721cf523101cd88eb07927c7e909281fb66` baseline was checked against GitHub: all 4,219 tracked blob hashes matched, with tree `20ec844dc688c1e13e7307cef4d4ae9d212bfca3`. The new implementation adds an owner page and read-only inventory; original routes and feature implementations remain present.

No current Qiniu catalog, existing provider key, provider account, production database or application route was read or changed. The preview contains synthetic identity and no live usage/cost evidence. Test-only fixture changes retain all six existing quota platforms and keep local fake SDK requests isolated; no production quota/provider/CAPTCHA behavior is removed.

## Passed checks

- Frontend final full suite: **341 files / 2,607 tests passed**, single worker. The first run found three stale five-platform assertions in the original quota defaults fixture; adding the existing `typesafe` platform to the expected test data resolved them without changing implementation.
- Complete frontend typecheck (`vue-tsc --noEmit` and `vue-tsc -b`) and production Vite build passed with Node heap 3,072 MiB and serialized compiler use. The earlier 1,536 MiB attempt exhausted Node heap; it was not an assertion failure.
- Whole frontend ESLint passed, without `--fix`.
- New owner inventory and route tests passed under Go 1.27.0, including foreign/admin-API-key identity denial, positive authenticated subject, explicit field projection, secret/error redaction, pagination and read-only route registration.
- Twenty top-level retained gateway protocol fixture tests passed through `tools/everplain/verify-model-contract.sh`: Responses and Chat Completions terminal usage, cache, tools/function arguments, returned model/service tier and interrupted streams. An initial service-package compile was killed under resource pressure; retry with serial compilation and bounded GC/memory passed.
- Three read-only preview tests passed: null provider cost evidence, unavailable usage and rejected writes.
- Three static source-overlay tests passed: commit-tagged fork build, port replacement/loopback binding and independent staging container names. These tests require PyYAML and do not run containers.
- Original Aliyun CAPTCHA fixture tests passed after adding their exact fake loopback host:port to test-local NO_PROXY, preserving existing exclusions. The SDK compares NO_PROXY by exact host:port. Before isolation, three fixtures timed out through the inherited egress proxy. Complete repository package retest passed; production verifier logic was untouched.

## Pending or unverified

- Complete backend unit aggregate exited **1**: 57 packages passed, but the original service package failed two parameter-matrix subtests because their public vendor hostnames could not resolve, and three plugin-host tests because Unix sockets are denied in this sandbox (`socket: operation not permitted`). Do not disable DNS/SSRF validation or plugin capabilities to make the candidate pass. Rerun on CI with the required capabilities. A proposed DNS-independent test fixture was not published because its recompile was interrupted by resource pressure. That statement records the local attempt only; the standard CI run above subsequently passed the complete unit/integration suites and golangci-lint.
- The full Go embed binary build did not pass. The first attempt ended without a usable binary; the retry reported `/tmp` exhaustion. A complete backend production binary is not claimed. Further large local Go checks are paused to avoid blocking other work; only regenerable compilation caches and exited build directories were removed, while source and check results were retained.
- The standard CI/security passes are linked above. The additional embedded-binary build must be checked on its own new source head. No permission expansion is introduced by this increment.
- Browser pixel review: my cloud browser returned `net::ERR_BLOCKED_BY_CLIENT` for the loopback-only preview. Source/component checks and a successful build do not verify screenshots, mobile layout or dark-mode pixels.
- Docker Compose rendered configuration and container startup were not verified. No staging instance was provisioned.
- Real Everplain-to-gateway-to-provider tests, credentials, canonical list-price verification, provider debit reconciliation and traffic cutover are outstanding. Actual upstream charges remain unavailable.

## No-provider-key runtime prerequisites

An actual management instance can start without any upstream model key: account inventory is empty and generation cannot be verified. Normal management mode still requires PostgreSQL, Redis, persistent application configuration and a securely initialized administrator. The source compose baseline uses PostgreSQL 18 and Redis 8 on internal ports 5432/6379, with the gateway on internal port 8080; the independent staging overlay replaces the exposed app binding with `127.0.0.1:8088` by default. It does not expose database/Redis ports or connect the current Everplain/Qiniu routes.

The source also supports a first-run setup wizard (`AUTO_SETUP` unset/false), which can serve the embedded frontend before database initialization. A setup-only listener is not a working administrator management instance. Keep first-run setup on loopback and complete operator-controlled initialization before any public reverse-proxy exposure.

This cloud workspace was checked and has no Docker, PostgreSQL or Redis executables and no complete gateway binary from the interrupted local build. It is not evidence that the user's Mac or an existing deployment host lacks these services. An operator must select an existing authorized host and verify its database, Redis, persistent storage and port/reverse-proxy state without copying secrets into the source or this report. No paid host, new credential, database instance or production route has been created here. Upstream-key migration and live model/cost checks remain separate from a keyless management deployment.

The downloadable preview runs the compiled production Vue assets through a separate read-only Python server. Start it with `python3 preview.py --owner-preview`; it does not grant a real administrator session or connect a model provider.
