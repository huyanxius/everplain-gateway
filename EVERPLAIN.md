# Everplain Gateway

A small, single-tenant, server-to-server adaptation of [Sub2API](https://github.com/Wei-Shaw/sub2api), pinned at `b8dece9000c68815a5b867ca5a1e6f236e173905` (v0.2.13).

Everplain keeps the end-user credit ledger. This service supplies model access through the existing Go/Gin/Ent + PostgreSQL + Redis architecture; the Vue frontend is embedded into the same server binary. No new runtime microservice is introduced.

- Accounts, provider adapters, model mappings, concurrency controls and usage persistence are retained.
- Upstream calls start **disabled**. AGY/Antigravity has no configured real credentials and no claimed readiness.
- Payments, redemption, subscriptions, affiliate/promotion entry points and public registration are blocked on the server as well as hidden in the UI.
- Simple mode is an operational **soft budget**, not a hard token/spend limit or tenant-isolation boundary. This instance must be an Everplain-only provider pool.
- No production deployment or paid model calls are part of this delivery.

## Start here

1. Read [the API and usage contract](docs/everplain/CONTRACT.md).
2. Build with [the development guide](docs/everplain/DEVELOPMENT.md), Go **1.27.0**, and pnpm.
3. Use [the lightweight configuration](deploy/everplain/config.example.yaml); do not inherit the upstream high-throughput example pools by accident.
4. Review [provider authorization and safe onboarding](docs/everplain/PROVIDER_SETUP.md); including the retained manual Antigravity OAuth flow, its side effects, and applicable terms.
5. Read [upstream provenance and launch/license checks](docs/everplain/UPSTREAM.md).

The original README files, LICENSE and attribution remain unchanged. The upstream README contains a separate commercial-authorization statement and provider-terms warnings. This fork does not claim that either issue has been cleared for a commercial launch.
