import { beforeEach, describe, expect, it, vi } from 'vitest'
import { getEverplainStatus } from '../everplain'
import { apiClient } from '../../client'
vi.mock('../../client', () => ({ apiClient: { get: vi.fn() } }))
const status = {
  contract_version: '2026-10-04', enabled: true,
  provider: { state: 'unconfigured', verification: 'not_performed', agy_adapter: 'antigravity_reserved' },
  upstream_enabled: false, ledger_owner: 'everplain', budget_enforcement: 'soft',
  requests_per_minute: 60, request_timeout_seconds: 120, usage_consistency: 'eventual', automatic_client_retries: false,
}
beforeEach(() => vi.clearAllMocks())
describe('gateway status contract', () => {
  it('uses only the existing admin-authenticated read endpoint', async () => {
    vi.mocked(apiClient.get).mockResolvedValue({ data: status })
    const controller = new AbortController()
    expect(await getEverplainStatus(controller.signal)).toEqual(status)
    expect(apiClient.get).toHaveBeenCalledWith('/admin/everplain/status', { signal: controller.signal })
  })
  it('keeps configured credentials explicitly unverified', async () => {
    const configured = { ...status, provider: { ...status.provider, state: 'configured_unverified' } }
    vi.mocked(apiClient.get).mockResolvedValue({ data: configured })
    expect((await getEverplainStatus()).provider.state).toBe('configured_unverified')
  })
  it.each([
    { contract_version: 'future' },
    { provider: { ...status.provider, state: 'ready' } },
    { provider: { ...status.provider, verification: 'passed' } },
    { automatic_client_retries: true },
    { budget_enforcement: 'hard' },
  ])('rejects an incompatible status rather than implying readiness: %j', async patch => {
    vi.mocked(apiClient.get).mockResolvedValue({ data: { ...status, ...patch } })
    await expect(getEverplainStatus()).rejects.toThrow('Unsupported gateway status contract')
  })
})
