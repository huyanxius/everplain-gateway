# Everplain gateway frontend

This is a small Vue 3 / TypeScript / Vite / Tailwind / Pinia adaptation of the pinned upstream frontend. The account editor, group/model routing, API key workflows, usage queries and settings continue to use their upstream API clients and state logic.

## Design source

The canonical Web `tokens.css` is copied byte-for-byte; its Git blob is verified by a test. `src/styles/everplain/SOURCE.md` records the exact original token and component source. `adapter.css` maps existing Vue classes to those tokens. Canvas charts use an auxiliary palette extracted from the same CSS, with a source-digest and value-parity test. There are no proprietary bundled fonts.

## Standalone profile

- Six admin navigation entries: overview, upstream accounts, model routing (upstream groups), API keys, usage, settings.
- Overview reads `/api/v1/admin/everplain/status` and existing dashboard counters. Unknown/incompatible/error results are explicit errors, never provider-ready states.
- AGY remains reserved. Credentials that the server reports as configured remain explicitly unverified. The UI does not perform provider verification or automatic generation probes.
- Payment, orders, subscriptions, redemption, promotion/referral, public registration, OAuth landing routes, and model plaza are closed before view loading and omitted from the route registry, so their lazy chunks are not shipped. Disabled routes are also excluded from prefetching. Stale public settings cannot reopen the closed consumer features.
- Backend authorization and feature closure are the enforcement boundary. The frontend profile is an additional presentation and routing restriction.
- Upstream optional commerce behavior remains in source for easier updates. Its retained upstream tests run without the standalone profile, while separate profile/route tests enforce the gateway closure.
- The overview discloses eventual usage consistency, soft budget enforcement, no automatic client retries, and no verified provider connection.

## Verification commands

Use the repository's pnpm workflow with its frozen lockfile:

```sh
pnpm install --frozen-lockfile
pnpm typecheck
pnpm test:run --maxWorkers=1 --minWorkers=1 --no-file-parallelism
pnpm build --outDir dist
pnpm lint:check
```

The explicit `dist` output keeps a standalone frontend check from modifying the backend's embedded assets. Integration builds can use the unchanged upstream default output path. The build's i18n gate uses one worker to avoid exhausting constrained CI/cloud environments.

## Acceptance boundaries

No real upstream account, key, paid model call, payment, or production deployment is used by these checks. Component tests cover initial loading, configuration states, repeated clicks, refresh after failure, and abort on navigation. Static source checks verify light/dark token selection, reduced-motion treatment, control typography, and mobile drawer keyboard support.

A real-browser visual review still needs an allowed local preview environment. In this cloud workspace a fresh headless Chromium launch failed on a restricted `socket()` syscall, and the supported cloud browser rejected the localhost preview with `ERR_BLOCKED_BY_CLIENT`. No external preview was published to work around that restriction. Thus pixel-level light/dark/narrow-screen parity, clipping, and native focus-ring rendering are not claimed verified.
