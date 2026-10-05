# Everplain frontend

The complete pinned Vue 3 / TypeScript / Vite / Tailwind / Pinia frontend is retained. All upstream routes and account/provider/model/key/usage/admin/user/payment/registration/settings flows remain available under upstream permissions and configuration. No extra lightweight mode or hardcoded feature closure is present.

The canonical Everplain Web tokens.css is copied byte-for-byte (blob e0a71a0197feab710e8540923cd0834fa32ea15b). SOURCE.md records provenance. adapter.css maps upstream presentation onto these tokens; chart canvas colors are derived from the same source. Branding and accessibility changes do not remove business handlers or navigation targets.

Visual acceptance is not inferred from token presence. The first actual browser screenshots exposed dark foregrounds mapped to surfaces, scoped table-header hardcoding, and a preview-only overlay banner. Corrections are made at the source, but require a fresh bundle and screenshot review before final visual acceptance.

The separately packaged preview server is read-only and contains no real credentials or usage. Production authentication and account setup are the unchanged upstream workflows; see docs/everplain/PROVIDER_SETUP.md before using them.
