# Unified gateway candidate

This additive candidate retains the complete pinned upstream distribution. It does not edit the current Qiniu model directory, rotate either upstream key, switch production routing, or replace Everplain's application ledger.

## Owner management entry

`/admin/everplain-gateway` is a new owner-only entry that links to the original account, group, pricing and personal-key editors. The owner administrator's default landing and logo/home links use this personal page; the original `/admin/dashboard` remains explicitly available as System dashboard. Existing admin/user pages and all upstream functions remain available. Other administrators retain their original landing.

The new `GET /api/v1/admin/everplain-gateway/upstreams` requires a validated administrator JWT for `huyanxius@gmail.com`, including a positive authenticated user ID. An admin API key is deliberately insufficient for this personal view. Its explicit projection exposes account IDs as `upstream-ID`, platform/type/configured status, scheduling flag, routing group IDs and safe mapped model IDs. It does not expose account names, credentials, notes, arbitrary extras, URLs, error messages, or invoke provider/billing probes. The alias identifies an existing account; it neither rotates nor creates its credential. More than two accounts can be managed without redesign.

Statistics use the existing `/usage/stats` and `/usage/dashboard/models` endpoints. Their backend binds user identity to the authenticated subject, rather than a selectable admin user ID. The new page never calls the system-wide admin dashboard. Other existing administrator pages retain their original authorization and statistical scope.

## Three cost meanings

1. `total_cost` / model `cost`: estimate at the configured gateway standard price. It is an official list-price estimate only after the exact model, input/output/cache tariff, currency, effective date and official source are checked.
2. `total_actual_cost` / model `actual_cost`: deduction in the gateway's existing internal ledger. Account allocations, including `account_cost`, remain internal configured accounting calculations.
3. Actual Qiniu/UniGate/provider charge: unavailable unless provider-issued debit/invoice evidence is reconciled to the request or bounded test interval. The new response intentionally uses `amount:null`, `currency:null`, `source:null`; absent evidence is never zero cost. No UniGate procurement rate or discount is inferred from a public vendor list price.

No new billing database or settlement path is added. Everplain's existing user-facing ledger and gateway operational accounting can coexist; neither proves a provider actually debited the same amount. Failed/cancelled/retried attempts may still consume provider tokens and must be reconciled separately.

## Everplain application protocol contract

Reviewed against Everplain `13e854890e48` source: `adapters/research_agent/pydantic_runner.py`, `adapters/model/metering.py`, `adapters/model/token_usage.py` and provider settings. Luna uses Responses; background summaries can use Chat Completions. Gemini/UniGate configurations remain independent. Integration must preserve both protocols, tools/structured output, and model-specific settings.

- Responses SSE must retain output/function-argument deltas and a single terminal response with numeric final usage. An error, truncated stream, missing terminal usage or contradictory counters must remain unknown/error; a synthetic success or usage estimate is not acceptable.
- Chat Completions streaming must retain tool deltas and the final usage snapshot requested through `stream_options.include_usage=true` (including an empty choices usage-only chunk). Terminal finish reasons and cancellation remain meaningful.
- Preserve the actual request/returned model, provider response ID and service tier. Preserve total input/output tokens and cache/reasoning subsets without double counting. Gateway and Everplain accounting may have different internal representations; raw provider evidence remains the authority.
- Existing gateway model mapping can rewrite the client-visible response model to its requested alias. For the first integration, use unchanged canonical model IDs and no nonidentity mapping. A nonidentity alias must not be treated as evidence of the provider-returned model. If aliases are required later, add and test an explicit opt-in preservation contract before routing billable traffic.
- Avoid protocol auto-detection/fallback and forced service-tier or reasoning rewrites for the initial contract test. Configure the original account's Responses support/mode and passthrough controls deliberately after confirming provider capabilities. A UI status or source fixture cannot verify those capabilities.

`tools/everplain/verify-model-contract.sh` runs existing upstream credential-free tests for Responses/Chat Completions, terminal usage, cache, service-tier observations, mapped response models, client tools, function arguments and interrupted streams. It is fixture regression coverage, not a live-provider or application-to-provider acceptance claim.

## Isolated source build and next gates

The original deployment compose file pulls `weishaw/sub2api:latest`; that is not this fork. `deploy/docker-compose.everplain-source.yml` overlays a source build, tagged with the verified source commit, separate container names and a loopback-only port. Combine both compose files under the new project name `everplain-gateway-staging`; this creates project-scoped volumes distinct from an existing project. Compose >= 2.24.4 is needed for the port override. Check the fully rendered compose configuration before starting it. No staging or production instance has been started by this change.

The original auto-setup provisions administrator/authentication state. The operator must supply and securely enter the required database, administrator, JWT and TOTP values and approve any persistent access; do not copy live secrets into source, archives, logs or chat. Use fresh isolated PostgreSQL/Redis. Credentials are never included in the candidate archive.

Release gates:

1. Final-head unit/integration/lint, full frontend tests/type/build and required CI must pass. Record any unrun or resource-blocked check accurately.
2. Review the owner page and original navigation in a real browser, light/dark and mobile, including refresh, failed requests, identity changes and interrupted navigation.
3. On the isolated instance, the owner securely configures the already-selected upstream credentials, models and limits. This is separate from current Qiniu model work and does not change the existing key or route.
4. Approve a bounded live contract test with explicit model/protocol, prompt data, budget and provider destination. Run Responses text/tools/structured output, Chat Completions usage and cancel/error cases through the real Everplain consumer, comparing request IDs, returned model/service tier and token/cache facts at both ledgers.
5. Reconcile the same interval against provider billing evidence; verify canonical model list prices independently. Only then authorize a separate traffic cutover with rollback to the untouched current route.

The retained upstream self-update entry remains available, but may download an upstream binary and overwrite fork-specific changes. Do not use it as the fork's deployment procedure.
