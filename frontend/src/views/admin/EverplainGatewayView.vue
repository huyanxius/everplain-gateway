<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-2xl font-semibold">{{ t('everplain.gateway.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.description') }}</p>
        </div>
        <button v-if="owner" type="button" class="btn btn-secondary" :disabled="loading" @click="load">
          {{ t('everplain.refresh') }}
        </button>
      </div>

      <div v-if="!owner" class="card p-6" role="alert">{{ t('everplain.gateway.ownerOnly') }}</div>
      <template v-else>
        <div class="card space-y-2 p-4 text-sm">
          <p>{{ t('everplain.gateway.scope') }}</p>
          <p class="text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.readiness') }}</p>
        </div>

        <div v-if="inventoryError" class="card p-4" role="alert">{{ t('everplain.gateway.inventoryError') }}</div>
        <div v-else-if="loading && !inventory" class="flex justify-center py-8" aria-live="polite"><LoadingSpinner /></div>

        <template v-if="inventory">
          <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <RouterLink v-for="link in links" :key="link.path" :to="link.path" class="card p-4 transition-colors hover:border-primary-500">
              <span class="font-medium">{{ t(link.title) }}</span>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t(link.description) }}</p>
            </RouterLink>
          </div>

          <section class="card p-4" :aria-label="t('everplain.gateway.ownUsage')">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <h2 class="text-lg font-medium">{{ t('everplain.gateway.ownUsage') }}</h2>
              <label class="flex items-center gap-2 text-sm">
                {{ t('everplain.gateway.period') }}
                <select v-model="days" class="input w-auto" :disabled="loading" @change="load">
                  <option :value="1">{{ t('everplain.gateway.today') }}</option>
                  <option :value="7">{{ t('everplain.gateway.week') }}</option>
                  <option :value="30">{{ t('everplain.gateway.month') }}</option>
                </select>
              </label>
            </div>
            <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ range.start_date }} → {{ range.end_date }} · {{ range.timezone }}</p>
            <p v-if="usageError" class="mt-3" role="alert">{{ t('everplain.gateway.usageError') }}</p>
            <div v-else-if="stats" class="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div><p class="text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.requests') }}</p><p class="mt-1 text-2xl font-semibold">{{ number(stats.total_requests) }}</p></div>
              <div><p class="text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.tokens') }}</p><p class="mt-1 text-2xl font-semibold">{{ number(stats.total_tokens) }}</p></div>
              <div><p class="text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.estimate') }}</p><p class="mt-1 text-2xl font-semibold">{{ usd(stats.total_cost) }}</p></div>
              <div><p class="text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.ledgerCharge') }}</p><p class="mt-1 text-2xl font-semibold">{{ usd(stats.total_actual_cost) }}</p></div>
            </div>
            <p class="mt-4 text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.estimateNote') }}</p>
            <div class="mt-3 rounded-lg border border-gray-200 p-3 dark:border-gray-700">
              <p class="font-medium">{{ t('everplain.gateway.providerCharge') }}: {{ t('everplain.gateway.unavailable') }}</p>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.providerChargeNote') }}</p>
            </div>
          </section>

          <section class="card overflow-hidden" :aria-label="t('everplain.gateway.upstreams')">
            <div class="flex flex-wrap items-center justify-between gap-3 p-4">
              <h2 class="text-lg font-medium">{{ t('everplain.gateway.upstreams') }} · {{ inventory.total }}</h2>
              <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('everplain.gateway.aliasNote') }}</span>
            </div>
            <p v-if="!inventory.items.length" class="p-4 text-sm">{{ t('everplain.emptyActivity') }}</p>
            <div v-else class="overflow-x-auto">
              <table class="w-full text-left text-sm">
                <thead class="bg-gray-50 text-gray-500 dark:bg-gray-800 dark:text-gray-400"><tr>
                  <th class="p-3">{{ t('everplain.gateway.alias') }}</th><th class="p-3">{{ t('everplain.gateway.platform') }}</th><th class="p-3">{{ t('everplain.gateway.status') }}</th><th class="p-3">{{ t('everplain.gateway.groups') }}</th><th class="p-3">{{ t('everplain.gateway.models') }}</th>
                </tr></thead>
                <tbody><tr v-for="upstream in inventory.items" :key="upstream.id" class="border-t border-gray-100 dark:border-gray-700">
                  <td class="p-3 font-medium">upstream-{{ upstream.id }}</td>
                  <td class="p-3">{{ upstream.platform }} · {{ upstream.type }}</td>
                  <td class="p-3">{{ upstream.status }} · {{ upstream.schedulable ? t('everplain.gateway.schedulable') : t('everplain.gateway.paused') }}</td>
                  <td class="p-3">{{ upstream.group_ids.join(', ') || '—' }}</td>
                  <td class="max-w-xs p-3 break-words">{{ upstream.mapped_models.join(', ') || t('everplain.gateway.noMapping') }}</td>
                </tr></tbody>
              </table>
            </div>
            <div v-if="inventory.total > inventory.page_size" class="flex items-center justify-end gap-3 p-4">
              <button type="button" class="btn btn-secondary" :disabled="loading || page <= 1" @click="changePage(-1)">{{ t('common.back') }}</button>
              <span>{{ page }} / {{ Math.ceil(inventory.total / inventory.page_size) }}</span>
              <button type="button" class="btn btn-secondary" :disabled="loading || page * inventory.page_size >= inventory.total" @click="changePage(1)">{{ t('common.next') }}</button>
            </div>
          </section>

          <section v-if="!usageError && models.length" class="card overflow-hidden" :aria-label="t('everplain.gateway.modelUsage')">
            <h2 class="p-4 text-lg font-medium">{{ t('everplain.gateway.modelUsage') }}</h2>
            <div class="overflow-x-auto"><table class="w-full text-left text-sm">
              <thead class="bg-gray-50 text-gray-500 dark:bg-gray-800 dark:text-gray-400"><tr><th class="p-3">{{ t('everplain.gateway.models') }}</th><th class="p-3">{{ t('everplain.requests') }}</th><th class="p-3">{{ t('everplain.gateway.tokens') }}</th><th class="p-3">{{ t('everplain.gateway.estimate') }}</th><th class="p-3">{{ t('everplain.gateway.ledgerCharge') }}</th></tr></thead>
              <tbody><tr v-for="model in models" :key="model.model" class="border-t border-gray-100 dark:border-gray-700"><td class="p-3">{{ model.model }}</td><td class="p-3">{{ number(model.requests) }}</td><td class="p-3">{{ number(model.total_tokens) }}</td><td class="p-3">{{ usd(model.cost) }}</td><td class="p-3">{{ usd(model.actual_cost) }}</td></tr></tbody>
            </table></div>
          </section>
        </template>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { getGatewayInventory, getGatewayOwnUsage, isGatewayOwner } from '@/api/admin/everplainGateway'
import type { GatewayInventory, GatewayUsageRange } from '@/api/admin/everplainGateway'
import type { ModelStat, UsageStatsResponse } from '@/types'

const { t } = useI18n()
const auth = useAuthStore()
const owner = computed(() => isGatewayOwner(auth.user))
const loading = ref(false)
const inventoryError = ref(false)
const usageError = ref(false)
const inventory = ref<GatewayInventory | null>(null)
const stats = ref<UsageStatsResponse | null>(null)
const models = ref<ModelStat[]>([])
const days = ref(7)
const page = ref(1)
let activeRequest: AbortController | undefined

function currentRange(): GatewayUsageRange {
  const end = new Date()
  const start = new Date(end)
  start.setDate(start.getDate() - days.value + 1)
  const date = (value: Date) => `${value.getFullYear()}-${String(value.getMonth() + 1).padStart(2, '0')}-${String(value.getDate()).padStart(2, '0')}`
  return { start_date: date(start), end_date: date(end), timezone: Intl.DateTimeFormat().resolvedOptions().timeZone }
}
const range = ref(currentRange())
const links = [
  { path: '/admin/accounts', title: 'everplain.upstreamAccounts', description: 'everplain.accountsDescription' },
  { path: '/admin/groups', title: 'everplain.modelRouting', description: 'everplain.routesDescription' },
  { path: '/admin/channels/pricing', title: 'everplain.gateway.pricing', description: 'everplain.gateway.pricingDescription' },
  { path: '/keys', title: 'everplain.keys', description: 'everplain.keysDescription' }
]
const number = (value: number) => Number.isFinite(value) ? value.toLocaleString() : '—'
const usd = (value: number) => Number.isFinite(value) ? `$${value.toFixed(6)}` : '—'

async function load() {
  activeRequest?.abort()
  const request = new AbortController()
  activeRequest = request
  inventory.value = null
  stats.value = null
  models.value = []
  if (!owner.value) { loading.value = false; return }
  loading.value = true
  inventoryError.value = false
  usageError.value = false
  range.value = currentRange()
  try {
    const result = await getGatewayInventory(page.value, request.signal)
    if (request.signal.aborted) return
    inventory.value = result
    try {
      const usage = await getGatewayOwnUsage(range.value, request.signal)
      if (request.signal.aborted) return
      stats.value = usage.stats
      models.value = usage.models
    } catch {
      if (!request.signal.aborted) usageError.value = true
    }
  } catch {
    if (!request.signal.aborted) inventoryError.value = true
  } finally {
    if (activeRequest === request) loading.value = false
  }
}
function changePage(delta: number) { page.value += delta; void load() }
watch(() => [auth.user?.id, auth.user?.email, auth.user?.role], () => { page.value = 1; void load() })
onMounted(load)
onUnmounted(() => activeRequest?.abort())
</script>
