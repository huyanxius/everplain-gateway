# Everplain Web design source

Retrieved from `huyanxius/everplain` via the authorized GitHub connector on 2026-10-04.

- Canonical tokens: `frontend/src/styles/tokens.css`, Git blob `e0a71a0197feab710e8540923cd0834fa32ea15b`.
- Component reference: `frontend/src/styles/components.css`, Git blob `d338cd00449cbe489615bfd3c4c81139a1a13780`.
- Source links: https://github.com/huyanxius/everplain/blob/main/frontend/src/styles/tokens.css and https://github.com/huyanxius/everplain/blob/main/frontend/src/styles/components.css.

`tokens.css` is copied verbatim from the canonical Web source. Do not edit it to redesign this distribution. Reconcile changes against the source first. The prior native Swift/Kotlin generated tokens were not used as the authoring source.

`adapter.css` maps the upstream Vue class names onto Everplain's semantic surfaces, ink accent, serif headings, 13px minimum metadata, pill controls, object-specific radii, focus ring, spacing, shadows, and motion. The Tailwind config keeps upstream utility names with token-backed palettes. Colored status/category utilities retain their semantics. `html.dark` selects `color-scheme: dark`, which makes the original `light-dark()` values honor the saved appearance choice.

The gateway overview is a new read-only composition of the original component conventions. Existing account, group/model routing, API key, usage, and settings logic stays in Vue and is not replaced with a mock implementation. The gateway profile closes commerce and public registration routes in addition to the backend boundary; it is not an authorization substitute.

Modern browser support for CSS `light-dark()` and `color-mix()` is required, matching the canonical Web token source. Fonts use the original system fallback stacks; no proprietary fonts are redistributed.

The default brand mark preserves the four paths from the previously verified Web `brand.svg` asset used for the native Everplain app. Only its fill is mapped to current token ink; no replacement logo geometry is invented. The wordmark is the existing product name in the source reading-font stack. The favicon uses the same geometry and canonical light/dark ink values.

`chart-palette.json` is auxiliary generated canvas data from the literal `light-dark()` pairs in `tokens.css`, not a separate authoring source. Its test verifies the canonical Git blob and every pair.
