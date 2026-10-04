<template>
  <AppLayout>
    <div class="everplain-overview">
      <section class="everplain-status card" :aria-busy="loading" aria-labelledby="gateway-status-title">
        <div class="everplain-section-head">
          <div>
            <h2 id="gateway-status-title">{{ t('everplain.gatewayStatus') }}</h2>
            <p>{{ t('everplain.gatewayStatusDescription') }}</p>
          </div>
          <button type="button" class="btn btn-secondary" :disabled="loading" @click="refresh">
            <Icon name="refresh" size="sm" />
            {{ t('everplain.refresh') }}
          </button>
        </div>

        <div v-if="loading && !status" class="everplain-status-loading" role="status">
          <LoadingSpinner size="sm" />
          <span>{{ t('common.loading') }}</span>
        </div>
        <div v-else-if="statusError" class="everplain-notice everplain-notice-danger" role="alert">
          <h3>{{ t('everplain.statusUnavailable') }}</h3>
          <p>{{ t('everplain.statusUnavailableDescription') }}</p>
        </div>
        <template v-else-if="status">
          <div v-if="!status.enabled" class="everplain-notice" role="status">{{ t('everplain.profileInactive') }}</div>
          <div class="everplain-provider">
            <div>
              <span class="everplain-meta">{{ t('everplain.agyReserved') }}</span>
              <h3>{{ t(status.provider.state === 'unconfigured' ? 'everplain.notConfigured' : 'everplain.configuredUnverified') }}</h3>
              <p>{{ t('everplain.providerDescription') }}</p>
            </div>
            <span class="everplain-state"><span aria-hidden="true"></span>{{ t(status.upstream_enabled ? 'everplain.upstreamEnabled' : 'everplain.upstreamDisabled') }}</span>
          </div>
          <div class="everplain-policy-grid">
            <div>
              <h3>{{ t('everplain.policy') }}</h3>
              <dl>
                <div><dt>{{ t('everplain.rateLimit') }}</dt><dd>{{ status.requests_per_minute }}</dd></div>
                <div><dt>{{ t('everplain.requestTimeout') }}</dt><dd>{{ status.request_timeout_seconds }} {{ t('everplain.seconds') }}</dd></div>
              </dl>
              <p>{{ t('everplain.noAutomaticRetries') }}</p>
            </div>
            <div>
              <h3>{{ t('everplain.accounting') }}</h3>
              <dl>
                <div><dt>{{ t('everplain.ledger') }}</dt><dd>{{ status.ledger_owner }}</dd></div>
                <div><dt>{{ t('everplain.softBudget') }}</dt><dd>{{ t('common.enabled') }}</dd></div>
              </dl>
              <p>{{ t('everplain.eventualUsage') }}</p>
            </div>
          </div>
          <p v-if="checkedAt" class="everplain-meta everplain-checked-at">{{ t('everplain.statusTime') }} · {{ checkedAt }}</p>
        </template>
      </section>

      <section aria-labelledby="gateway-activity-title">
        <h2 id="gateway-activity-title" class="everplain-section-title">{{ t('everplain.activity') }}</h2>
        <div v-if="activityError" class="everplain-notice" role="status">{{ t('everplain.activityUnavailable') }}</div>
        <div v-else class="everplain-metrics" :aria-busy="loading">
          <div v-for="metric in metrics" :key="metric.label" class="card">
            <span class="everplain-meta">{{ metric.label }}</span>
            <strong>{{ metric.value === undefined ? '—' : new Intl.NumberFormat(locale).format(metric.value) }}</strong>
          </div>
        </div>
        <p v-if="stats?.total_accounts === 0" class="everplain-empty-hint">{{ t('everplain.emptyActivity') }}</p>
      </section>

      <section aria-labelledby="gateway-workspace-title">
        <h2 id="gateway-workspace-title" class="everplain-section-title">{{ t('everplain.workspace') }}</h2>
        <div class="everplain-workspace-grid">
          <router-link v-for="entry in workspace" :key="entry.path" :to="entry.path" class="card everplain-workspace-link">
            <Icon :name="entry.icon" size="md" aria-hidden="true" />
            <h3>{{ t(entry.title) }}</h3>
            <p>{{ t(entry.description) }}</p>
            <span class="everplain-open" aria-hidden="true">↗</span>
          </router-link>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getEverplainStatus, type EverplainStatus } from '@/api/admin/everplain'
import { getStats } from '@/api/admin/dashboard'
import type { DashboardStats } from '@/types'

const { t, locale } = useI18n()
const loading = ref(false)
const statusError = ref(false)
const activityError = ref(false)
const status = ref<EverplainStatus | null>(null)
const stats = ref<DashboardStats | null>(null)
const checkedAt = ref('')
let controller: AbortController | undefined
let disposed = false
const workspace = [
  { path: '/admin/accounts', title: 'everplain.upstreamAccounts', description: 'everplain.accountsDescription', icon: 'server' },
  { path: '/admin/groups', title: 'everplain.modelRouting', description: 'everplain.routesDescription', icon: 'grid' },
  { path: '/keys', title: 'nav.apiKeys', description: 'everplain.keysDescription', icon: 'key' },
  { path: '/admin/usage', title: 'nav.usage', description: 'everplain.usageDescription', icon: 'chart' },
] as const
const metrics = computed(() => [
  { label: t('everplain.accounts'), value: stats.value?.total_accounts },
  { label: t('everplain.keys'), value: stats.value?.total_api_keys },
  { label: t('everplain.requests'), value: stats.value?.total_requests },
])

async function refresh() {
  if (loading.value) return
  loading.value = true
  statusError.value = false
  activityError.value = false
  controller = new AbortController()
  const [statusResult, statsResult] = await Promise.allSettled([
    getEverplainStatus(controller.signal), getStats(),
  ])
  if (disposed) return
  if (statusResult.status === 'fulfilled') {
    status.value = statusResult.value
    checkedAt.value = new Date().toLocaleTimeString(locale.value, { hour: '2-digit', minute: '2-digit' })
  } else {
    status.value = null
    statusError.value = true
  }
  if (statsResult.status === 'fulfilled') stats.value = statsResult.value
  else { stats.value = null; activityError.value = true }
  loading.value = false
}
onMounted(refresh)
onBeforeUnmount(() => { disposed = true; controller?.abort() })
</script>
