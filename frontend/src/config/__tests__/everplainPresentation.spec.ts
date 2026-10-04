import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
const read = (path: string) => readFileSync(resolve(process.cwd(), 'src', path), 'utf8')
const css = read('styles/everplain/adapter.css')
const sidebar = read('components/layout/AppSidebar.vue')
describe('Everplain presentation and accessibility contract', () => {
  it('overrides the original root scheme with enough specificity for an explicit light preference', () => {
    expect(css).toContain(':root { color-scheme: light; }')
    expect(css).toContain(':root.dark { color-scheme: dark; }')
  })
  it('uses original semantic font, control shape, focus, and reduced-motion tokens', () => {
    for (const token of ['--qx-font-reading', '--qx-font-ui', '--qx-radius-pill', '--qx-radius-card', '--qx-motion-base']) expect(css).toContain(token)
    expect(css).toContain(':focus-visible')
    expect(css).toContain('(prefers-reduced-motion: reduce)')
    expect(css).toContain('animation-iteration-count: 1 !important')
    expect(css).toContain('(max-width: 479px)')
  })
  it('keeps the closed mobile sidebar inert and supports Escape, focus wrap, and focus return', () => {
    expect(sidebar).toContain(':inert="isMobileViewport && !mobileOpen"')
    expect(sidebar).toContain("event.key === 'Escape'")
    expect(sidebar).toContain("event.key !== 'Tab'")
    expect(sidebar).toContain('previousFocus?.focus()')
    expect(sidebar).toContain('document.removeEventListener')
  })
})
