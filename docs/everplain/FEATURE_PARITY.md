# Full baseline and optional simplification candidates

Reviewed against upstream `b8dece9000c68815a5b867ca5a1e6f236e173905` on 2026-10-04. The implemented candidate is the full baseline. A lighter alternative is a proposal only: no additional feature is disabled or deleted without the operator's item-by-item choice.

## Source evidence

A file-hash comparison of the current working tree with the fixed upstream base found:

- Router registry: byte-identical; all 66 path declarations retained.
- `frontend/src/views/`: 195 upstream files (including tests), none missing; only AccountsView differs, by an informational notice.
- `frontend/src/api/`: all 85 upstream files retained, none changed.
- `backend/internal/handler/`: all 360 upstream files retained, none changed.
- `backend/internal/service/`: all 1,375 upstream files retained, none changed.
- `backend/internal/server/routes/`: all 23 upstream files retained, none changed.
- Backend config and server router restored to upstream. The only backend source change is the tested-in-prior-checkpoint setup DSN fix and its regression test.
- The custom Everplain route guard, billing/usage wrapper, mandatory simple mode and hardcoded six-entry navigation have been removed.
- VersionBadge remains visible and retains version details, update checks, update and rollback actions. Its earlier accidental hiding was corrected. Upstream binary self-update may overwrite this fork; the action remains available and is not executed by this delivery.

Counts describe source preservation, not executed end-to-end acceptance. They include helper/test files under each directory and must not be described as numbers of UI pages or endpoints.

## Functional matrix

| Capability | Source evidence | Current source state | What remains unverified |
|---|---|---|---|
| Admin/user routes, navigation and page actions | router/index.ts; AppSidebar.vue; AppHeader.vue | Original route registry/menu branches restored; brand/a11y changes retained | Full browser navigation, role permutations |
| Provider accounts, OAuth and token refresh | service/*oauth_service.go; token_refresh_service.go; account components | Original implementation retained, no credential supplied | Actual provider grant/refresh; operator-controlled only |
| API compatibility, model routing and scheduling | routes/gateway.go; gateway_service.go; openai_gateway_service.go; openai_account_scheduler.go | Original aliases, protocol handlers and schedulers retained | Every SDK/provider combination |
| Concurrency, health, retries and quota protection | concurrency_service.go; model_rate_limit.go; openai_apikey_health_breaker.go; antigravity_gateway_retry.go; middleware/api_key_auth.go | Original mechanisms retained | Full failure/load tests on this candidate |
| Usage, cost, credits and audit/security | gateway_usage_billing.go; usage_service.go; audit_log_service.go; account_credentials_redact.go | Original standard/simple semantics retained | End-to-end accounting/reconciliation/security suites |
| Images, edits, async/batch image tasks | routes/gateway.go; original image pages/services | No text-only filter; handlers/pages retained | Per-provider availability and credentials |
| Video, Grok TTS/STT/realtime, Responses WebSocket | routes/gateway.go and provider services | Original endpoints/aliases retained | Real media/realtime operation |
| Payment, orders, sales packages and subscriptions | routes/payment.go; payment_order.go; payment_order_expiry_service.go | Original APIs, page routes/settings restored | Real purchase/webhook/refund; not authorized here |
| Public registration and end-user social login | routes/auth.go; auth_service.go; original auth views | Original settings and callbacks restored | External login/registration; not performed |
| Promotions, affiliates, redemption and model plaza | routes/admin.go; routes/user.go; promo_service.go; affiliate_service.go; redeem_service.go | Original code and menus restored | Transactional integrations and permission combinations |
| Operational settings/plugins/proxies/monitoring | original SettingsView.vue, admin route files and services | Settings loaders/tabs and admin navigation restored | Every toggle and background-job lifecycle |
| Version/update/rollback | VersionBadge.vue | Entry and original behavior retained | Fork-safe update target selection; no update run |

Provider OAuth is not the same thing as an end user's social login. Everplain owning its own customer ledger is not a reason to remove gateway usage records, cost accounting, quotas or security protections. Simple mode changes accounting and some group-isolation semantics; forcing it is not an ordinary performance optimization.

## Comparison with the earlier checkpoint

The old light frontend blocked payment/orders/redemption/subscriptions/affiliates/model plaza/registration/auth/batch-image and overwrote public OAuth flags. Its backend admitted only three POST text endpoints plus a few GET routes, excluding images, video, audio, WebSocket, aliases and update/rollback. That scope was broader than removing retail features. Those additional restrictions are absent from the full candidate. Restoring source does not automatically configure live providers, enable every upstream business switch or authorize external operations.

## Lighter alternative: three candidates for operator choice

There is no absolute must-delete feature and no measured CPU/RAM saving for these candidates. Keep the code by default; configuration-only choices are easier to reverse than exclusion from builds or deletion.

1. **Payment/orders/sales packages.** The second retail checkout may duplicate Everplain's existing business flow. Its real maintenance includes webhook verification, refunds, idempotency, fulfillment and reconciliation; background expiry/reconciliation scanning exists on a roughly 60-second cycle. Merely hiding payment UI does not stop these jobs. After confirmation, disable collection settings while retaining code. Removing packaging would require explicit DI/job separation and tests. Configuration reversal is low effort; hard deletion is high risk.
2. **Public signup and unneeded end-user social logins.** Server-only consumption may not need public account creation. The maintenance burden includes verification email, abuse prevention and identity linking/callbacks. Preserve administrator login, 2FA/passkeys, recovery and provider OAuth. The original DingTalk path can differ from the global registration switch, so hiding a registration page is not proof that signup is closed. Confirm exact switches and backend behavior first.
3. **Promotional codes/affiliate retail credits/redemption sales entry points.** These may duplicate Everplain's customer marketing. The maintenance includes incentive rules, affiliate-to-balance transfers and consistent redemption fulfillment. Preserve operational credit tools and schema. Do not delete by a `redeem` string match: the admin account code also has a provider-specific Claude quota-reset redemption action. Configuration is relatively reversible; deleting accounting/database relationships is not.

Multimodal and protocol aliases are not proposed for deletion: future image understanding, voice and realtime work may use them, and removal can break SDK compatibility. Any approved reduction needs a separate change set, tests and rollback plan; this document authorizes none.

## Acceptance boundary

Independent review covered source differences only. Current SFC compile/lint/token checks and historical test results are recorded separately in VERIFICATION.md. Full current build/type/suite checks, functional browser coverage, real-provider access and full-version resource measurements remain outstanding. No production deployment, credentials, model calls or purchases were performed.
