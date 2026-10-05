import { beforeEach, describe, expect, it, vi } from 'vitest'
import { getAdminLandingPath, getGatewayInventory, getGatewayOwnUsage, isGatewayOwner } from '../admin/everplainGateway'

const { get } = vi.hoisted(() => ({ get: vi.fn() }))
vi.mock('../client', () => ({ apiClient: { get } }))

describe('Everplain gateway read-only API contract', () => {
  beforeEach(() => get.mockReset())
  it('allows only the exact administrator owner', () => {
    expect(isGatewayOwner(null)).toBe(false)
    expect(isGatewayOwner({ role: 'user', email: 'huyanxius@gmail.com' })).toBe(false)
    expect(isGatewayOwner({ role: 'admin', email: 'other@gmail.com' })).toBe(false)
    expect(isGatewayOwner({ role: 'admin', email: 'huyanxius@gmail.com.evil.test' })).toBe(false)
    expect(isGatewayOwner({ role: 'admin', email: 'huyanxius@Gmail.com' })).toBe(true)
  })
  it('uses only the safe inventory endpoint with bounded pagination', async () => {
    const signal = new AbortController().signal
    get.mockResolvedValue({ data: { items: [] } })
    await getGatewayInventory(2, signal)
    expect(get).toHaveBeenCalledTimes(1)
    expect(get).toHaveBeenCalledWith('/admin/everplain-gateway/upstreams', { params: { page: 2, page_size: 25 }, signal })
  })
  it('makes the personal gateway the owner landing while preserving the original dashboard', () => {
    expect(getAdminLandingPath({ role: 'admin', email: 'huyanxius@gmail.com' })).toBe('/admin/everplain-gateway')
    expect(getAdminLandingPath({ role: 'admin', email: 'other@gmail.com' })).toBe('/admin/dashboard')
    expect(getAdminLandingPath(null)).toBe('/admin/dashboard')
  })
  it('does not forward user selectors to own-usage endpoints', async () => {
    get.mockResolvedValue({ data: { models: [] } })
    const range = { start_date: '2026-10-01', end_date: '2026-10-05', timezone: 'UTC', user_id: 999, account_id: 8 }
    await getGatewayOwnUsage(range)
    expect(get).toHaveBeenNthCalledWith(1, '/usage/stats', { params: { start_date: range.start_date, end_date: range.end_date, timezone: 'UTC' }, signal: undefined })
    expect(get).toHaveBeenNthCalledWith(2, '/usage/dashboard/models', { params: { start_date: range.start_date, end_date: range.end_date, timezone: 'UTC' }, signal: undefined })
    expect(get.mock.calls.some(([path]) => path.includes('/admin/dashboard'))).toBe(false)
  })
})
