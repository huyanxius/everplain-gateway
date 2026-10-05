import { beforeEach, describe, expect, it, vi } from 'vitest'
import { reactive } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import EverplainGatewayView from '../EverplainGatewayView.vue'

const mocks = vi.hoisted(() => ({
  inventory: vi.fn(), usage: vi.fn(), auth: null as { user: { email: string; role: string } | null } | null
}))
vi.mock('@/stores/auth', () => ({ useAuthStore: () => mocks.auth }))
vi.mock('vue-i18n', async (importOriginal) => ({
  ...await importOriginal<typeof import('vue-i18n')>(),
  useI18n: () => ({ t: (key: string) => key })
}))
vi.mock('@/api/admin/everplainGateway', async (importOriginal) => {
  const original = await importOriginal<typeof import('@/api/admin/everplainGateway')>()
  return { ...original, getGatewayInventory: mocks.inventory, getGatewayOwnUsage: mocks.usage }
})

const inventory = () => ({
  items: [{ id: 3, alias: 'ignored-account-name', platform: 'openai', type: 'apikey', status: 'active', schedulable: true, group_ids: [2], mapped_models: ['luna'] }],
  count: 1, total: 1, page: 1, page_size: 25,
  provider_charge: { status: 'unavailable', amount: null, currency: null, source: null, reason: 'unavailable' }
})
const usage = () => ({
  stats: { total_requests: 5, total_tokens: 50, total_cost: 0.004, total_actual_cost: 0.003 },
  models: [{ model: 'luna', requests: 5, total_tokens: 50, cost: 0.004, actual_cost: 0.003 }]
})
const render = () => mount(EverplainGatewayView, { global: { stubs: { AppLayout: { template: '<main><slot /></main>' }, LoadingSpinner: true, RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' } } } })

describe('Owner gateway view', () => {
  beforeEach(() => {
    mocks.inventory.mockReset().mockResolvedValue(inventory())
    mocks.usage.mockReset().mockResolvedValue(usage())
    mocks.auth = reactive({ user: { role: 'admin', email: 'huyanxius@Gmail.com' } })
  })
  it('shows only an upstream alias and separates ledger from unknown provider cost', async () => {
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('upstream-3')
    expect(wrapper.text()).not.toContain('ignored-account-name')
    expect(wrapper.text()).toContain('$0.004000')
    expect(wrapper.text()).toContain('$0.003000')
    expect(wrapper.text()).toContain('everplain.gateway.providerCharge: everplain.gateway.unavailable')
    expect(wrapper.find('a[href="/admin/accounts"]').exists()).toBe(true)
    expect(wrapper.find('a[href="/admin/groups"]').exists()).toBe(true)
    wrapper.unmount()
  })
  it('loads neither inventory nor usage for another admin', async () => {
    mocks.auth!.user!.email = 'other@gmail.com'
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('everplain.gateway.ownerOnly')
    expect(mocks.inventory).not.toHaveBeenCalled()
    expect(mocks.usage).not.toHaveBeenCalled()
    wrapper.unmount()
  })
  it('does not load usage after inventory authorization fails or display raw server errors', async () => {
    mocks.inventory.mockRejectedValue(new Error('private upstream detail'))
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('everplain.gateway.inventoryError')
    expect(wrapper.text()).not.toContain('private upstream detail')
    expect(mocks.usage).not.toHaveBeenCalled()
    wrapper.unmount()
  })
  it('keeps inventory usable while hiding failed usage totals', async () => {
    mocks.usage.mockRejectedValue(new Error('private database error'))
    const wrapper = render()
    await flushPromises()
    expect(wrapper.text()).toContain('upstream-3')
    expect(wrapper.text()).toContain('everplain.gateway.usageError')
    expect(wrapper.text()).not.toContain('$0.000000')
    expect(wrapper.text()).not.toContain('private database error')
    wrapper.unmount()
  })
  it('aborts pending work on unmount and does not request usage afterward', async () => {
    let resolve: (value: ReturnType<typeof inventory>) => void = () => {}
    mocks.inventory.mockReturnValue(new Promise((done) => { resolve = done }))
    const wrapper = render()
    const signal = mocks.inventory.mock.calls[0][1] as AbortSignal
    wrapper.unmount()
    expect(signal.aborted).toBe(true)
    resolve(inventory())
    await flushPromises()
    expect(mocks.usage).not.toHaveBeenCalled()
  })
  it('clears owner data and aborts requests when the signed-in identity changes', async () => {
    const wrapper = render()
    await flushPromises()
    const signal = mocks.inventory.mock.calls[0][1] as AbortSignal
    mocks.auth!.user = { role: 'admin', email: 'other@gmail.com' }
    await flushPromises()
    expect(signal.aborted).toBe(true)
    expect(wrapper.text()).not.toContain('upstream-3')
    expect(wrapper.text()).toContain('everplain.gateway.ownerOnly')
    wrapper.unmount()
  })
  it('aborts and clears data when the owner email changes on the same user object', async () => {
    const wrapper = render()
    await flushPromises()
    const signal = mocks.inventory.mock.calls[0][1] as AbortSignal
    mocks.auth!.user!.email = 'other@gmail.com'
    await flushPromises()
    expect(signal.aborted).toBe(true)
    expect(wrapper.text()).not.toContain('upstream-3')
    expect(wrapper.text()).toContain('everplain.gateway.ownerOnly')
    wrapper.unmount()
  })
})
