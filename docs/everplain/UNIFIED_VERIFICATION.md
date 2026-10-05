# Unified management verification

2026-10-05. Source candidate, not production acceptance.

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

- Complete backend unit run is still being checked; do not substitute focused/package passes for final aggregate success. Integration and golangci-lint are not yet verified on this source candidate.
- No GitHub Actions pass is claimed. No workflow enablement or permission expansion was performed.
- Browser pixel review: my cloud browser returned `net::ERR_BLOCKED_BY_CLIENT` for the loopback-only preview. Source/component checks and a successful build do not verify screenshots, mobile layout or dark-mode pixels.
- Docker Compose rendered configuration and container startup were not verified. No staging instance was provisioned.
- Real Everplain-to-gateway-to-provider tests, credentials, canonical list-price verification, provider debit reconciliation and traffic cutover are outstanding. Actual upstream charges remain unavailable.

The downloadable preview runs the compiled production Vue assets through a separate read-only Python server. Start it with `python3 preview.py --owner-preview`; it does not grant a real administrator session or connect a model provider.
