import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import EverplainOverviewView from '../EverplainOverviewView.vue'
import { getEverplainStatus } from '@/api/admin/everplain'
import { getStats } from '@/api/admin/dashboard'
vi.mock('@/components/layout/AppLayout.vue', () => ({ default: { template: '<main><slot /></main>' } }))
vi.mock('@/api/admin/everplain', () => ({ getEverplainStatus: vi.fn() }))
vi.mock('@/api/admin/dashboard', () => ({ getStats: vi.fn() }))
vi.mock('vue-i18n', async (importOriginal) => ({ ...await importOriginal<typeof import('vue-i18n')>(), useI18n: () => ({ t: (key: string) => key, locale: { value: 'en' } }) }))
const status = {
  contract_version: '2026-10-04', enabled: true,
  provider: { state: 'unconfigured', verification: 'not_performed', agy_adapter: 'antigravity_reserved' },
  upstream_enabled: false, ledger_owner: 'everplain', budget_enforcement: 'soft',
  requests_per_minute: 60, request_timeout_seconds: 120, usage_consistency: 'eventual', automatic_client_retries: false,
} as const
function render() {
  return mount(EverplainOverviewView, {
    global: { stubs: { AppLayout: { template: '<main><slot /></main>' }, Icon: true, LoadingSpinner: { template: '<span role="status">Loading</span>' }, RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' } } },
  })
}
beforeEach(() => {
  vi.resetAllMocks()
  vi.mocked(getEverplainStatus).mockResolvedValue(status)
  vi.mocked(getStats).mockResolvedValue({ total_accounts: 0, total_api_keys: 0, total_requests: 0 } as Awaited<ReturnType<typeof getStats>>)
})
describe('Everplain gateway overview', () => {
  it('shows an unconfigured state and the soft/eventual accounting caveat', async () => {
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('everplain.notConfigured')
    expect(wrapper.text()).toContain('everplain.upstreamDisabled')
    expect(wrapper.text()).toContain('everplain.eventualUsage')
    expect(wrapper.text()).toContain('everplain.emptyActivity')
    expect(wrapper.findAll('.everplain-workspace-link')).toHaveLength(4)
    wrapper.unmount()
  })
  it('never labels configured credentials verified or healthy', async () => {
    vi.mocked(getEverplainStatus).mockResolvedValue({ ...status, provider: { ...status.provider, state: 'configured_unverified' } })
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('everplain.configuredUnverified')
    expect(wrapper.text()).not.toContain('everplain.ready')
    wrapper.unmount()
  })
  it('exposes loading and suppresses repeated refreshes', async () => {
    let resolve!: (value: typeof status) => void
    vi.mocked(getEverplainStatus).mockReturnValue(new Promise(r => { resolve = r }))
    const wrapper = render()
    await flushPromises()
    expect(wrapper.find('[aria-busy="true"]').exists()).toBe(true)
    await wrapper.get('button').trigger('click')
    await wrapper.get('button').trigger('click')
    expect(getEverplainStatus).toHaveBeenCalledTimes(1)
    resolve(status)
    await flushPromises()
    expect(wrapper.get('button').attributes('disabled')).toBeUndefined()
    wrapper.unmount()
  })
  it('reports status and metrics failures honestly and permits retry', async () => {
    vi.mocked(getEverplainStatus).mockRejectedValueOnce(new Error('offline'))
    vi.mocked(getStats).mockRejectedValueOnce(new Error('offline'))
    const wrapper = render()
    await flushPromises()
    expect(wrapper.get('[role="alert"]').text()).toContain('everplain.statusUnavailable')
    expect(wrapper.text()).toContain('everplain.activityUnavailable')
    expect(wrapper.text()).not.toContain('everplain.notConfigured')
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(wrapper.find('[role="alert"]').exists()).toBe(false)
    expect(wrapper.text()).toContain('everplain.notConfigured')
    wrapper.unmount()
  })
  it('cancels its status read when navigating away', async () => {
    vi.mocked(getEverplainStatus).mockReturnValue(new Promise(() => {}))
    const wrapper = render()
    const signal = vi.mocked(getEverplainStatus).mock.calls[0][0]
    wrapper.unmount()
    expect(signal?.aborted).toBe(true)
  })
})
