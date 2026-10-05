"""Static source-build overlay contract; not a Docker runtime/config test."""
from pathlib import Path
import unittest

import yaml


class OverrideList(list):
    pass


class ComposeLoader(yaml.SafeLoader):
    pass


ComposeLoader.add_constructor('!override', lambda loader, node: OverrideList(loader.construct_sequence(node)))


class EverplainSourceOverlayTest(unittest.TestCase):
    def setUp(self):
        path = Path(__file__).resolve().parents[1] / 'docker-compose.everplain-source.yml'
        self.services = yaml.load(path.read_text(), Loader=ComposeLoader)['services']

    def test_fork_build_is_commit_tagged_and_does_not_pull_upstream_latest(self):
        app = self.services['sub2api']
        self.assertEqual(app['build']['context'], '..')
        self.assertEqual(app['build']['dockerfile'], 'Dockerfile')
        self.assertIn('EVERPLAIN_SOURCE_COMMIT:?', app['image'])
        self.assertNotIn('weishaw/', app['image'])
        self.assertIn('EVERPLAIN_SOURCE_COMMIT:?', app['build']['args']['COMMIT'])

    def test_port_replaces_inherited_public_binding(self):
        ports = self.services['sub2api']['ports']
        self.assertIsInstance(ports, OverrideList)
        self.assertEqual(len(ports), 1)
        self.assertTrue(ports[0].startswith('127.0.0.1:'))
        self.assertIn('EVERPLAIN_STAGING_PORT:-8088', ports[0])

    def test_staging_containers_are_distinct_and_no_secret_is_embedded(self):
        names = [service['container_name'] for service in self.services.values()]
        self.assertEqual(len(names), len(set(names)))
        self.assertTrue(all(name.startswith('everplain-gateway-staging') for name in names))
        self.assertEqual(self.services['sub2api']['environment'], ['ADMIN_EMAIL=huyanxius@gmail.com'])


if __name__ == '__main__':
    unittest.main()
