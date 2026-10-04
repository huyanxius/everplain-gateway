# Everplain Gateway

An independent, branded distribution of [Sub2API](https://github.com/Wei-Shaw/sub2api), pinned at `b8dece9000c68815a5b867ca5a1e6f236e173905` (v0.2.13).

## Scope

All original upstream features, routes and modes are retained. Everplain changes presentation using its canonical Web tokens, brand mark and accessible navigation. The earlier text-only/light profile has been removed; there is no Everplain feature-disable switch, hardcoded six-item menu, commerce block or replacement billing contract. Upstream configuration and permissions continue to govern feature availability.

The Go/Gin/Ent + PostgreSQL + Redis architecture and Vue frontend are unchanged. Provider code, account actions, routing, quotas, billing, payment, registration/OAuth, user/admin pages, images/audio/video, WebSocket and operational settings remain the upstream implementations. No actual provider credentials are included and no live account, model request, purchase or production deployment has been performed.

## Read first

- [Build, preview and verification](docs/everplain/DEVELOPMENT.md)
- [Manual provider setup, OAuth side effects and terms](docs/everplain/PROVIDER_SETUP.md)
- [Functional preservation evidence](docs/everplain/FEATURE_PARITY.md)
- [Upstream provenance and licensing](docs/everplain/UPSTREAM.md)

The original upstream README, LICENSE and attribution remain unchanged. Provider terms and upstream's separate commercial-use statement require review before commercial launch; this repository is not legal clearance.
