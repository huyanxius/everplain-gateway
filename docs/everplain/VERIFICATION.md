# Verification ledger — full-feature candidate, 2026-10-04

This distinguishes current checks from the older restricted-profile checkpoint. No live account, provider request, payment or deployment was used.

## Passed on current source

- Canonical Everplain tokens.css exact Git blob: `e0a71a0197feab710e8540923cd0834fa32ea15b`; same current main source verified through GitHub.
- Vue SFC parse, script compilation and template compilation: all 10 changed components.
- ESLint: all 21 changed JavaScript/TypeScript/Vue files.
- Everplain presentation/a11y/source contract: 4 tests. Includes dark foreground and table token mapping regression checks.
- Everplain locale key equality: 48 matching English/Chinese keys.
- Git diff whitespace/error check.
- Preview HTTP smoke on the first-look package: real asset SPA bootstrap, visible demo banner, synthetic unconfigured status, usage unavailable (503), writes rejected (405). The later full-source preview wrapper is Python syntax checked; it still needs a new full frontend build and browser review.

## Blocked / not passed

- Full current frontend typecheck, full test suite and production build: processes were killed under cloud memory pressure (exit 137), including serialized bounded-memory attempts. They are not marked passing.
- Current backend setup test: compiler was killed while compiling the upstream service dependency. The DSN fix previously passed in the older checkpoint, but that is not a current full-build pass.
- Fresh full-feature browser screenshots: not available until the full source builds. The earlier Mac screenshots were of the six-menu checkpoint and exposed dark contrast/table/banner defects; source fixes are not equivalent to rendered acceptance.
- Full upstream integration/security/lint suites and full-version idle/load resources: not rerun successfully.
- CI: the public repository's Actions page and API showed zero runs/checks. Workflow files exist and CI triggers on push/pull_request without a draft exclusion. The operator-side browser subsequently confirmed “Workflows aren’t being run on this forked repository.” Fork workflows require activation before the existing push/PR checks can run. No permissions were expanded by this source change. The frontend CI job now adds the full suite, production asset build and a local-preview artifact, without deployment.

## Historical only

The previous restricted-profile UI passed 341 test files / 2647 tests plus typecheck/build. Synthetic provider load measurements and its 239–253 MiB peak RSS concerned that prior profile, not this full-feature candidate. Do not use these numbers as current full-version sizing or production guarantees.

The source remains a draft candidate until missing checks and real-browser review are completed. No production launch is authorized by this ledger.
