import importlib.util
import json
from pathlib import Path
import threading
import unittest
from urllib.error import HTTPError
from urllib.request import Request, urlopen

spec = importlib.util.spec_from_file_location('gateway_preview', Path(__file__).with_name('preview.py'))
preview = importlib.util.module_from_spec(spec)
spec.loader.exec_module(preview)


class PreviewContractTest(unittest.TestCase):
    def setUp(self):
        self.server = preview.ThreadingHTTPServer(('127.0.0.1', 0), preview.Preview)
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()
        self.base = f'http://127.0.0.1:{self.server.server_port}'

    def tearDown(self):
        self.server.shutdown()
        self.server.server_close()
        self.thread.join(timeout=2)

    def test_inventory_has_no_account_or_cost_evidence(self):
        with urlopen(self.base + '/api/v1/admin/everplain-gateway/upstreams') as response:
            data = json.load(response)['data']
            self.assertEqual(response.headers['Cache-Control'], 'no-store')
        self.assertEqual(data['items'], [])
        self.assertIsNone(data['provider_charge']['amount'])
        self.assertIsNone(data['provider_charge']['currency'])
        self.assertIsNone(data['provider_charge']['source'])

    def test_usage_is_unavailable_and_writes_are_rejected(self):
        for path, method, status in [('/api/v1/usage/stats', 'GET', 503),
                ('/api/v1/usage/dashboard/models', 'GET', 503),
                ('/api/v1/admin/accounts', 'POST', 405),
                ('/api/v1/admin/accounts/1', 'DELETE', 405)]:
            with self.subTest(path=path, method=method):
                with self.assertRaises(HTTPError) as failure:
                    urlopen(Request(self.base + path, method=method))
                self.assertEqual(failure.exception.code, status)

    def test_owner_preview_session_remains_explicitly_synthetic(self):
        self.assertEqual(preview.USER['id'], 0)
        self.assertNotEqual(preview.USER['email'], 'huyanxius@gmail.com')
        self.assertIn('NONSECRET-UI-PREVIEW-ONLY', preview.make_bootstrap(preview.USER))


if __name__ == '__main__':
    unittest.main()
