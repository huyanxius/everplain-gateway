import { apiClient } from '../client'
import type { ModelStat, UsageStatsResponse } from '@/types'

export interface GatewayUpstream {
  id: number
  alias: string
  platform: string
  type: string
  status: string
  schedulable: boolean
  group_ids: number[]
  mapped_models: string[]
}

export interface GatewayInventory {
  items: GatewayUpstream[]
  count: number
  total: number
  page: number
  page_size: number
  provider_charge: {
    status: 'unavailable'
    amount: null
    currency: null
    source: null
    reason: string
  }
}

export interface GatewayUsageRange {
  start_date: string
  end_date: string
  timezone: string
}

export function isGatewayOwner(user: { role: string; email: string } | null): boolean {
  return user?.role === 'admin' && user.email.trim().toLowerCase() === 'huyanxius@gmail.com'
}

export function getAdminLandingPath(user: { role: string; email: string } | null): string {
  return isGatewayOwner(user) ? '/admin/everplain-gateway' : '/admin/dashboard'
}

export async function getGatewayInventory(page = 1, signal?: AbortSignal): Promise<GatewayInventory> {
  const { data } = await apiClient.get<GatewayInventory>('/admin/everplain-gateway/upstreams', {
    params: { page, page_size: 25 }, signal
  })
  return data
}

// These existing user endpoints bind UserID to the authenticated subject on the server.
// Never forward an admin-selected user_id or use the system-wide dashboard here.
export async function getGatewayOwnUsage(range: GatewayUsageRange, signal?: AbortSignal) {
  const params = { start_date: range.start_date, end_date: range.end_date, timezone: range.timezone }
  const [stats, models] = await Promise.all([
    apiClient.get<UsageStatsResponse>('/usage/stats', { params, signal }),
    apiClient.get<{ models: ModelStat[] }>('/usage/dashboard/models', { params, signal })
  ])
  return { stats: stats.data, models: models.data.models }
}
