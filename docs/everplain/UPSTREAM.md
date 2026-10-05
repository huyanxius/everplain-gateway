# Upstream provenance and license notes

- Upstream: https://github.com/Wei-Shaw/sub2api
- Pinned base: `b8dece9000c68815a5b867ca5a1e6f236e173905`
- Commit: `chore: sync VERSION to 0.2.13 [skip ci]`, 2026-10-02T12:16:57Z.
- Required Go: **1.27.0**, from `backend/go.mod` and upstream CI. No guessed downgrade.
- Original `LICENSE`, `README.md`, `README_CN.md`, `README_JA.md`, and `CLA.md` are retained unchanged.
- Go module/import paths remain `github.com/Wei-Shaw/sub2api` to minimize fork churn.
- Ent/PostgreSQL, Redis, Gin, Vue/Vite, existing providers, scheduler, model mappings and usage persistence are retained. There is no additional runtime microservice.

## Practical release boundary

The upstream LICENSE is GNU LGPL v3. The README separately states that its developers have not authorized commercial operations based on the project, and warns about provider terms of service. This fork preserves both statements and does not claim commercial approval from upstream or from any model provider. Before commercial launch, clarify the README's commercial statement with upstream or qualified counsel, and verify the applicable provider account/API terms, especially any subscription-to-API resale use. This is an unresolved launch check, not a legal clearance.

If distributing an executable/container, preserve notices and make the corresponding modified source and applicable license information available in a way meeting the license's requirements. Exact obligations for a particular combined distribution need legal review; a public repository alone is not a blanket compliance guarantee. No production deployment or paid model call was performed by this work.

## Fork differences

See `docs/everplain/CONTRACT.md` and `docs/everplain/DEVELOPMENT.md`. Fork additions are deliberately small boundary/configuration changes; upstream architecture and migrations are not rewritten. UI adaptations stay in the existing frontend and are embedded into the same backend binary.

## Current feature-preservation scope

The earlier Everplain text-only/light profile has been removed at the operator's request. The full upstream route registry, admin/user menus, payment/registration and provider protocol handlers are preserved. In particular VersionBadge's version details, update and rollback actions remain accessible. Invoking upstream binary self-update can replace this fork with an upstream build; review the target before using it. This is a warning, not a disabled feature, and no update or rollback was performed.
