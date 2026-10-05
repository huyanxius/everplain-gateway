import { resolve } from 'node:path'
import { readFileSync } from 'node:fs'
import { createHash } from 'node:crypto'
import palette from '@/styles/everplain/chart-palette.json'
import { describe, expect, it } from 'vitest'
import { everplainChartColor } from '../useEverplainChartTheme'
describe('canvas palette follows canonical Web tokens', () => {
  it('matches the canonical Web source byte for byte and every generated chart pair', () => {
    const source = readFileSync(resolve(process.cwd(), 'src/styles/everplain/tokens.css'), 'utf8')
    const digest = createHash('sha1').update(`blob ${Buffer.byteLength(source)}\0${source}`).digest('hex')
    expect(digest).toBe(palette.source_blob)
    for (const [name, pair] of Object.entries(palette.colors)) {
      expect(source).toContain(`--qx-color-${name}: light-dark(${pair[0]}, ${pair[1]});`)
    }
  })
  it('uses exact light and dark category colors', () => {
    expect(everplainChartColor('data-blue')).toBe('#43639a')
    expect(everplainChartColor('data-blue', true)).toBe('#95a5c0')
    expect(everplainChartColor('ink-soft', true)).toBe('#c9ccc6')
  })
  it('fails visibly if a requested source token is removed', () => {
    expect(() => everplainChartColor('invented')).toThrow('Missing Everplain chart color token')
  })
})
