# Everplain server gateway contract

Contract version: `2026-10-04`.

## Ownership and trust boundary

Everplain's application backend is the sole end-user identity and credit ledger. It holds a gateway service key in its server-side secret storage. Never embed this key or any provider credential in the browser, macOS/iOS app, or public repository. This gateway is a dedicated **single-tenant provider pool for Everplain**, not an end-user wallet or multi-tenant billing service.

The profile requires upstream `run_mode: simple`. In this upstream version simple mode bypasses billing/quota checks **and some group isolation** and can schedule the whole eligible provider pool. Groups/model mappings remain useful operational configuration, but are not a tenant security boundary in this profile. Use separate instances if a future customer needs isolated accounts. API-key authentication, user/key disabled checks, allowed model validation and provider scheduling remain upstream code.

The profile has soft operational budget semantics. Its Redis RPM cap limits requests, not tokens or spend, and can admit a burst at a window boundary. Existing USD cost values are telemetry/estimates only. A generation may complete after client cancellation, usage records arrive asynchronously, and a failing write can lose telemetry. Do not advertise strict spending limits or exact billing guarantees. Everplain must reserve, settle and reconcile its own end-user credits under its own transaction/idempotency policy. This fork adds no credit debits.

## Default state

`everplain.enabled=true`, `run_mode=simple`, `everplain.upstream_enabled=false`. Requests to generate content return HTTP 503 `provider_not_configured` after service-key authentication. No default account, real credential or live provider call is introduced. AGY is an adapter placeholder pointing at the existing Antigravity implementation; account OAuth, privacy and onboarding operations have not been performed.

An administrator can read `GET /api/v1/admin/everplain/status` with the existing admin JWT or upstream admin authentication. Response uses the panel envelope `{code:0,data:{...}}`. `provider.state` is `unconfigured` while the upstream switch is off; when switched on it is `configured_unverified`. The switch is an operator acknowledgement, **not a successful credential/model readiness probe**. `provider.verification` remains `not_performed`; no synthetic success or green health badge is returned for an unverified account.

## Supported server-side routes

Use `Authorization: Bearer <service-key>` over TLS when outside localhost.

- `GET /v1/models` and `GET /v1/models/:model`: existing upstream model catalog/mapping behavior.
- `POST /v1/messages`: Anthropic-compatible request, response and SSE.
- `POST /v1/chat/completions`: OpenAI-compatible request, response and SSE.
- `POST /v1/responses`: OpenAI Responses request, response and SSE.
- `GET /v1/usage`: existing operational summary, whose historical labels may refer to quotas. These labels do not turn simple-mode budgets into hard limits.
- `GET /v1/everplain/usage?page=1&page_size=100`: normalized key-scoped usage events.

Generation still goes through existing authentication, group-model allowlists, model mapping, concurrency control, provider scheduler, and usage writing. The Everplain layer adds a distributed per-service-key RPM cap and a context deadline. It does **not** add another generation retry loop. The lightweight example reduces upstream account-switch attempts, but upstream-specific retry/failover behavior still exists.

Legacy root/provider-specific aliases, images, audio, video, web search, batch and WebSocket endpoints are unavailable in this text-only profile, preventing an alias from bypassing admission controls. Account/provider admin APIs remain available for explicitly authorized future setup. Payment/webhook/admin-payment, subscriptions, redemption, promotions, affiliate, model-plaza and public registration/OAuth endpoints are blocked server-side regardless of hidden UI links. Upstream binary self-update/rollback endpoints are blocked so they cannot overwrite the fork with an unrelated release. Login/2FA/refresh/logout and ordinary admin operations remain.

## Errors and retries

Before an SSE response starts, all gateway HTTP errors have a stable envelope:

```json
{"error":{"code":"provider_not_configured","type":"gateway_error","message":"No upstream has been enabled. Configure and verify an authorized provider first.","request_id":"opaque-gateway-correlation-id","retryable":false}}
```

`X-Client-Request-ID` is the gateway-generated correlation header; `X-Everplain-Contract-Version` identifies this contract. Provider error bodies are discarded rather than exposing tokens, prompt content or account details. Successful responses and SSE bytes are unmodified upstream protocol. Once SSE has started, errors remain provider-protocol terminal events or a closed stream; an HTTP JSON envelope cannot replace an already started stream. Consumers must handle both outcomes.

- 400/422: `invalid_request`, no automatic retry.
- 401/403: `authentication_failed` / `access_denied`, fix authorization.
- 404: `not_found` / `feature_disabled`.
- 413: `request_too_large`.
- 429: `rate_limited`; respect `Retry-After` for the Redis admission limit.
- 503 before admission: `provider_not_configured` (no retry), or `limiter_unavailable` / `provider_unavailable` (temporarily retryable).
- 499/504: cancellation / deadline. Generation outcome can be unknown.
- 500/502: internal / upstream error. Generation outcome can be unknown.

`retryable=true` describes temporary availability, not an idempotency guarantee. The API does not implement generation replay deduplication. Disable SDK automatic retries by default; only retry when the application knows no generation was admitted or accepts duplicate provider charges. Never retry after tokens have streamed merely because a connection closed. No change was made to upstream account retry decisions.

## Usage telemetry

The normalized endpoint returns:

```json
{"contract_version":"2026-10-04","ledger_owner":"everplain","consistency":"eventual","items":[{"usage_id":"123","request_id":"provider-request-id","session_id":"opaque-correlation","model":"requested-model","upstream_model":"mapped-model","input_tokens":10,"output_tokens":5,"cache_creation_tokens":0,"cache_read_tokens":0,"cost_usd":0.001,"actual_cost_usd":0.001,"stream":false,"created_at":"2026-10-04T00:00:00Z"}],"page":1,"page_size":100,"total":1}
```

`usage_id` is a string row ID for deduplication within this gateway database. `request_id` is the persisted upstream usage row identifier. When it has the upstream `client:` prefix, the additional `client_request_id` field exposes the suffix for direct matching to `X-Client-Request-ID`; historical or exceptional rows can omit that field. If correlation is needed, send an opaque, non-personal `X-Session-Id` under the existing upstream session contract; it may also participate in affinity. Avoid PII and prompts in any identifier.

The endpoint queries only the authenticated service key. No account/user IDs, credentials, prompt content, IP address or credits are exported. Each token category is explicit; don't add cache fields twice to provider-native counts. USD values follow upstream pricing and multiplier semantics and are not provider invoices. Failed reads return 503 `usage_unavailable`, never a fabricated empty usage list.

Pagination is newest-first and not a snapshot. New rows can shift pages; polling clients should overlap pages and deduplicate by `usage_id`. This initial contract is telemetry, not a lossless billing event stream or an exactly-once ledger. Client cancellation is propagated to the upstream request context, but cancellation does not prove the provider consumed zero tokens.
