import { computed, ref } from 'vue'
import { useMutationObserver } from '@vueuse/core'
import palette from '@/styles/everplain/chart-palette.json'

/** Canvas cannot reliably consume CSS variable expressions. Resolve only the
 * canonical literal light/dark pairs, so there is no second authored palette.
 */
export function everplainChartColor(name: string, dark = false): string {
  const pair = palette.colors[name as keyof typeof palette.colors]
  if (!pair) throw new Error(`Missing Everplain chart color token: ${name}`)
  return pair[dark ? 1 : 0]
}
export function useEverplainChartTheme() {
  const dark = ref(document.documentElement.classList.contains('dark'))
  useMutationObserver(document.documentElement, () => {
    dark.value = document.documentElement.classList.contains('dark')
  }, { attributes: true, attributeFilter: ['class'] })
  return computed(() => ({
    text: everplainChartColor('ink-soft', dark.value),
    grid: everplainChartColor('rule', dark.value),
    muted: everplainChartColor('faint', dark.value),
    series: ['blue', 'green', 'amber', 'rose', 'violet', 'teal', 'rust', 'olive'].map(name => everplainChartColor(`data-${name}`, dark.value)),
  }))
}
