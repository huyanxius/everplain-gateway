# Read-only Vue preview

The preview bundle includes compiled Vue assets and this standard-library Python server. No Node dependencies, PostgreSQL, Redis, provider credentials, or paid API are needed to view it. It runs only on loopback and is never a deployable gateway.

With Python 3 available, unzip the bundle, run `python3 preview.py` (Windows: `py preview.py`), then open http://127.0.0.1:4181/admin/dashboard. Stop with Ctrl+C. The banner sits in normal document flow and labels it as a UI preview. The synthetic administrator uses the original standard mode; example feature switches expose the full interface but do not connect payment, registration or provider services. The actual production Vue components render; the server supplies only clearly artificial identity, unconfigured status, and empty account lists. Usage and other unavailable APIs return errors rather than fabricated statistics. All writes return 405 and nothing is saved. Do not enter credentials.

The regular production app and its auth/API code are unchanged by preview mode. The bootstrap is injected only by this separate local preview server. Browser screenshot/pixel QA must be performed separately; a successful bundle build is not visual verification.

To inspect the new personal gateway entry, run `python3 preview.py --owner-preview` and open `/admin/everplain-gateway`. This explicitly synthetic owner session still has user ID zero and a nonsecret preview marker; it is not an authenticated production administrator. The upstream inventory is empty, actual provider charges are null, and usage remains unavailable. No existing Qiniu directory/key/route is read or modified.
