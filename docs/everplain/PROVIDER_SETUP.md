# Provider access and safe onboarding

Reviewed 2026-10-04. This is an engineering readiness note, not legal clearance.

## Antigravity is not a subscription-to-API entitlement

Google's [Antigravity terms](https://antigravity.google/terms), section 6, prohibit accessing the consumer service using third-party software with Antigravity OAuth and warn of suspension or termination. An account owner's permission does not change Google's service terms. This is a material service-terms risk for personal subscription OAuth in this gateway; the operator must decide whether to proceed.

The retained upstream implementation uses `cloudcode-pa.googleapis.com/v1internal:*` and Google's OAuth flow (`backend/internal/pkg/antigravity/oauth.go` and `client.go`). Its presence in open source does not make it a supported public Google model API. It is not the new official Antigravity SDK. The fork has not configured or verified this adapter, and upstream requests remain off by default.

Google now documents an [official Antigravity SDK](https://antigravity.google/docs/sdk/overview/) using a Gemini API key, or a Google Cloud project/location and authorized Cloud credentials. The consumer terms separately exclude certain Enterprise arrangements from their scope; that does not automatically approve this retained adapter or subscription-quota reuse. Review the actual product and organizational terms before implementation. An SDK-based agent runtime is a different integration from this gateway's model completion endpoints and is not implemented here.

## Supported next step

Choose a documented provider API whose terms permit the intended server-side use, such as Gemini API, Anthropic API, or OpenAI API. Existing upstream provider code is retained, but each selected API/model/protocol mapping still needs explicit integration verification; none is claimed connected by this delivery.

1. The operator selects the provider, project, model, intended use and spending limit. A paid account or subscription is not created automatically.
2. Create a dedicated least-privilege API credential in the provider's official console. Do not send passwords, cookies, refresh tokens, API keys or service-account secrets through chat, screenshots, repository commits or logs.
3. The operator enters the credential directly into the intended secure server-side secret store/admin handoff. Creating credentials or persistent access needs separate action-time approval. Do not copy credentials from the user's installed IDE or browser profile.
4. Verify the provider's allowed region, data use/retention, exact pricing and rate limits. Do not rotate accounts or retry to evade limits. Respect Retry-After.
5. On an isolated test instance, authorize a narrowly bounded live test. Check the selected request protocol, model mapping, cancellation, streaming terminal event, error redaction and persisted usage. Confirm provider-side usage and costs too.
6. Only after that test and a separate deployment decision may an operator enable upstream requests. The existing status `configured_unverified` is an acknowledgement of configuration, never proof of a passed provider check.

## Costs and limits

[Antigravity plans](https://antigravity.google/docs/plans) describe subscription quotas, weekly limits and optional AI-credit overages. These are not this gateway's API capacity. Limits can change and are not a guaranteed tokens-per-second allowance.

[Gemini API billing](https://ai.google.dev/gemini-api/docs/billing) is separately metered by project/billing account. Free-tier eligibility, paid tiers and available models vary. Do not promise that an AI Pro/Ultra subscription pays for API traffic. [Model pricing](https://ai.google.dev/gemini-api/docs/pricing) and the project's current rate-limit screen are the sources to check immediately before a paid test.

The gateway RPM limiter is not a monetary hard cap. Usage telemetry is eventually consistent; cancellation and retry can leave generation outcomes unknown. Everplain owns the end-user ledger. Never advertise exact billing guarantees based on this gateway's soft operational budgets.

## Current acceptance state

- No live credentials, upstream generation, account OAuth grant or paid API call has been made.
- No production deployment is authorized or performed.
- Synthetic fixture load measurements do not establish production provider capacity.
- Frontend uses the canonical Everplain Web tokens. Browser pixel QA remains blocked by the cloud browser's localhost access policy; do not represent source/component tests as screenshots.

## Manual path in the retained upstream UI

Open `/admin/accounts` → Add account → Antigravity → OAuth → Next → Generate authorization link. Open the generated Google page in your own browser and complete Google's consent step yourself if you choose to proceed. The callback is fixed in this pinned source to `http://localhost:8085/callback`; this is the browser operator's localhost, not the gateway domain. The UI accepts the complete callback URL or its code. State and PKCE must match the original server session, which expires after 30 minutes. Paste the callback only into the intended gateway panel, never into chat.

The actual endpoints are `POST /api/v1/admin/antigravity/oauth/auth-url`, then `POST /api/v1/admin/antigravity/oauth/exchange-code`. The latter exchanges the code, reads userinfo, calls LoadCodeAssist, may call OnboardUser, and invokes setUserSettings through the upstream privacy helper. After a successful exchange, the UI automatically creates the account and saves credentials. The refresh-token import is a separate path and also performs external requests. This delivery preserves that behavior and does not introduce a validation-only authorization mode.

`everplain.upstream_enabled: false` blocks public generation routes, but does not block admin OAuth, account tests, scheduled probes, privacy operations or background refresh. Therefore merely withholding a model prompt is not the same as no external requests or no account changes. Using the read-only preview bundle makes no Google requests; its write operations are disabled and it cannot authorize or create an account. An actual gateway instance requires the separate build, PostgreSQL and Redis setup described in DEVELOPMENT.md.

The operator chooses whether to use a Google/Antigravity account and completes all account steps themselves. No live authorization, token extraction, provider test or privacy change was performed as part of this delivery. Provider terms remain relevant even when no model generation is performed.
