#!/usr/bin/env python3
"""Read-only, loopback-only preview of the real Vue production build.

No provider calls, database, real auth, credentials, or usage records. Never deploy.
"""
import argparse
import json
from functools import partial
from http.server import SimpleHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import urlsplit

USER = {'id': 0, 'email': 'preview@example.invalid', 'username': '界面预览',
        'role': 'admin', 'status': 'active', 'balance': 0, 'concurrency': 0,
        'run_mode': 'standard'}
STATUS = {'contract_version': '2026-10-04', 'enabled': True,
          'provider': {'state': 'unconfigured', 'verification': 'not_performed',
                       'agy_adapter': 'antigravity_reserved'},
          'upstream_enabled': False, 'ledger_owner': 'everplain',
          'budget_enforcement': 'soft', 'requests_per_minute': 60,
          'request_timeout_seconds': 120, 'usage_consistency': 'eventual',
          'automatic_client_retries': False}
def make_bootstrap(user):
    return '''<script>
localStorage.setItem('auth_token','NONSECRET-UI-PREVIEW-ONLY');
localStorage.setItem('auth_user',%s);
localStorage.setItem('sub2api_locale','zh');
localStorage.removeItem('refresh_token');
localStorage.removeItem('token_expires_at');
</script>''' % json.dumps(json.dumps(user, ensure_ascii=False), ensure_ascii=False)
BOOTSTRAP = make_bootstrap(USER)
BANNER = """<style>@media(min-width:1024px){#everplain-preview-banner{margin-left:calc(var(--qx-outline-width) + var(--qx-space-4))!important}}</style><aside id="everplain-preview-banner" role="status" style="position:relative;margin:var(--qx-space-4);padding:var(--qx-space-4) var(--qx-space-6);border:1px solid var(--qx-color-rule);border-radius:var(--qx-radius-card);background:var(--qx-color-surface);color:var(--qx-color-ink);font-family:var(--qx-font-ui);font-size:var(--qx-text-control);line-height:var(--qx-text-control--line-height);text-align:center">只读界面预览 · 使用真实 Vue 页面和 Everplain tokens<br>状态为演示配置，未连接账号；用量不可用，操作不会保存。请勿输入任何密钥。</aside>"""


class Preview(SimpleHTTPRequestHandler):
    def log_message(self, fmt, *args):
        # Do not log arbitrary request data or headers.
        pass

    def json_response(self, status, data):
        body = json.dumps(data, ensure_ascii=False).encode()
        self.send_response(status)
        self.send_header('Content-Type', 'application/json; charset=utf-8')
        self.send_header('Cache-Control', 'no-store')
        self.send_header('Content-Length', str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_POST(self):
        self.json_response(405, {'code': 405, 'message': '只读预览：操作不会保存，请勿输入凭据。'})
    do_PUT = do_PATCH = do_DELETE = do_POST

    def do_GET(self):
        path = urlsplit(self.path).path
        if path == '/api/v1/auth/me':
            return self.json_response(200, {'code': 0, 'data': USER})
        if path == '/api/v1/admin/everplain/status':
            return self.json_response(200, {'code': 0, 'data': STATUS})
        if path == '/api/v1/admin/everplain-gateway/upstreams':
            return self.json_response(200, {'code': 0, 'data': {
                'items': [], 'count': 0, 'total': 0, 'page': 1, 'page_size': 25,
                'provider_charge': {'status': 'unavailable', 'amount': None,
                    'currency': None, 'source': None,
                    'reason': 'Read-only UI preview; no provider or billing evidence.'}}})
        if path in ('/api/v1/settings/public', '/api/v1/public/settings'):
            return self.json_response(200, {'code': 0, 'data': {
                'site_name': 'Everplain Gateway', 'registration_enabled': True,
                'payment_enabled': True, 'subscription_enabled': True,
                'affiliate_enabled': True, 'model_plaza_enabled': True,
                'backend_mode_enabled': False}})
        if path in ('/api/v1/setup/status', '/setup/status'):
            return self.json_response(200, {'code': 0, 'data': {'needs_setup': False}})
        if path in ('/api/v1/admin/accounts', '/api/v1/admin/groups', '/api/v1/keys'):
            return self.json_response(200, {'code': 0, 'data': {
                'items': [], 'total': 0, 'page': 1, 'page_size': 20, 'pages': 0}})
        if path in ('/api/v1/admin/groups/all', '/api/v1/groups/available'):
            return self.json_response(200, {'code': 0, 'data': []})
        if path.startswith(('/api/', '/v1/', '/v1beta/')):
            return self.json_response(503, {'code': 503, 'message': '界面预览未连接后端；无真实数据。'})
        candidate = Path(self.translate_path(path))
        if not candidate.is_file() or path == '/index.html':
            body = (Path(self.directory) / 'index.html').read_text()
            body = body.replace('<head>', '<head>' + BOOTSTRAP).replace('<div id="app"></div>', BANNER + '<div id="app"></div>')
            data = body.encode()
            self.send_response(200)
            self.send_header('Content-Type', 'text/html; charset=utf-8')
            self.send_header('Cache-Control', 'no-store')
            self.send_header('Content-Length', str(len(data)))
            self.end_headers()
            self.wfile.write(data)
            return
        return super().do_GET()

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--dist', type=Path, default=Path(__file__).resolve().parent / 'dist')
    parser.add_argument('--port', type=int, default=4181)
    parser.add_argument('--owner-preview', action='store_true',
                        help='Display the owner UI using an explicitly synthetic id=0 session.')
    args = parser.parse_args()
    if args.owner_preview:
        USER = {**USER, 'email': 'huyanxius@gmail.com'}
        BOOTSTRAP = make_bootstrap(USER)
    if not (args.dist / 'index.html').is_file():
        parser.error('Expected an existing Vue dist/index.html; use the supplied preview bundle.')
    server = ThreadingHTTPServer(('127.0.0.1', args.port), partial(Preview, directory=str(args.dist.resolve())))
    landing = '/admin/everplain-gateway' if args.owner_preview else '/admin/dashboard'
    print(f'Read-only preview: http://127.0.0.1:{args.port}{landing}', flush=True)
    print('No live credentials or usage; do not deploy this preview server.', flush=True)
    server.serve_forever()
