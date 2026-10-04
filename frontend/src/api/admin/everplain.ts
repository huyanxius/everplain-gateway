import { apiClient } from '../client'

/** Versioned, read-only contract. Configuration does not establish provider health. */
export interface EverplainStatus {
  contract_version: string
  enabled: boolean
  provider: {
    state: 'unconfigured' | 'configured_unverified'
    verification: 'not_performed'
    agy_adapter: 'antigravity_reserved'
  }
  upstream_enabled: boolean
  ledger_owner: string
  budget_enforcement: 'soft'
  requests_per_minute: number
  request_timeout_seconds: number
  usage_consistency: 'eventual'
  automatic_client_retries: false
}

export async function getEverplainStatus(signal?: AbortSignal): Promise<EverplainStatus> {
  const { data } = await apiClient.get<EverplainStatus>('/admin/everplain/status', { signal })
  if (data.contract_version !== '2026-10-04' ||
      !['unconfigured', 'configured_unverified'].includes(data.provider?.state) ||
      data.provider?.verification !== 'not_performed' ||
      data.provider?.agy_adapter !== 'antigravity_reserved' ||
      typeof data.enabled !== 'boolean' || typeof data.upstream_enabled !== 'boolean' ||
      !Number.isFinite(data.requests_per_minute) || data.requests_per_minute < 0 ||
      !Number.isFinite(data.request_timeout_seconds) || data.request_timeout_seconds <= 0 ||
      data.budget_enforcement !== 'soft' || data.usage_consistency !== 'eventual' ||
      data.automatic_client_retries !== false) {
    throw new Error('Unsupported gateway status contract')
  }
  return data
}
