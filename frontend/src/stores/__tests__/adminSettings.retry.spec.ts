import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAdminSettingsStore } from '../adminSettings'

const profile = vi.hoisted(() => ({ enabled: false }))
vi.mock('@/config/everplain', () => ({ get EVERPLAIN_GATEWAY_PROFILE() { return profile.enabled } }))

const mocks = vi.hoisted(() => ({ getSettings: vi.fn(), getConfig: vi.fn() }))
vi.mock('@/api', () => ({ adminAPI: { settings: { getSettings: mocks.getSettings }, payment: { getConfig: mocks.getConfig } } }))
beforeEach(() => {
  profile.enabled = false
  vi.clearAllMocks()
  localStorage.clear()
  setActivePinia(createPinia())
  mocks.getSettings.mockResolvedValue({ ops_monitoring_enabled: true, custom_menu_items: [{ id: 'custom', title: 'Custom' }] })
  mocks.getConfig.mockResolvedValue({ data: { enabled: true } })
  vi.spyOn(console, 'error').mockImplementation(() => {})
})
afterEach(() => { vi.restoreAllMocks(); localStorage.clear() })

describe('admin settings fetch retry', () => {
  it.each(['settings', 'payment'])('retries a failed initial %s request without forcing a refresh', async (source) => {
    const request = source === 'settings' ? mocks.getSettings : mocks.getConfig
    request.mockRejectedValueOnce(new Error('Temporarily unavailable'))
    const store = useAdminSettingsStore()
    await store.fetch()
    expect(store.loaded).toBe(false)
    expect(store.loading).toBe(false)
    await store.fetch()
    expect(mocks.getSettings).toHaveBeenCalledTimes(2)
    expect(mocks.getConfig).toHaveBeenCalledTimes(2)
    expect(store.loaded).toBe(true)
    expect(store.paymentEnabled).toBe(true)
    expect(store.customMenuItems).toEqual([{ id: 'custom', title: 'Custom' }])
  })

  it('keeps cached values visible when the initial request fails', async () => {
    localStorage.setItem('ops_monitoring_enabled_cached', 'false')
    localStorage.setItem('payment_enabled_cached', 'true')
    mocks.getSettings.mockRejectedValueOnce(new Error('Offline'))
    const store = useAdminSettingsStore()
    await store.fetch()
    expect(store.opsMonitoringEnabled).toBe(false)
    expect(store.paymentEnabled).toBe(true)
    expect(localStorage.getItem('payment_enabled_cached')).toBe('true')
  })

  it('still reuses successfully loaded settings', async () => {
    const store = useAdminSettingsStore()
    await store.fetch()
    await store.fetch()
    expect(mocks.getSettings).toHaveBeenCalledTimes(1)
    expect(mocks.getConfig).toHaveBeenCalledTimes(1)
  })
})


describe('standalone gateway admin settings', () => {
  it('never reads the closed payment endpoint, even if it would fail', async () => {
    profile.enabled = true
    mocks.getConfig.mockRejectedValue(new Error('Payment is disabled'))
    const store = useAdminSettingsStore()
    await store.fetch()
    expect(store.loaded).toBe(true)
    expect(store.paymentEnabled).toBe(false)
    expect(mocks.getConfig).not.toHaveBeenCalled()
    expect(store.customMenuItems).toEqual([{ id: 'custom', title: 'Custom' }])
  })
  it('still retries a genuine settings failure under the gateway profile', async () => {
    profile.enabled = true
    mocks.getSettings.mockRejectedValueOnce(new Error('Offline'))
    const store = useAdminSettingsStore()
    await store.fetch()
    expect(store.loaded).toBe(false)
    await store.fetch()
    expect(store.loaded).toBe(true)
    expect(mocks.getSettings).toHaveBeenCalledTimes(2)
    expect(mocks.getConfig).not.toHaveBeenCalled()
  })
})
