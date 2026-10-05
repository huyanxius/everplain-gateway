import { describe, expect, it } from 'vitest'
import { resolveCompletedSetupRedirectPath } from '../setupRedirect'

describe('Completed setup keeps owner landing personal', () => {
  it('keeps login and non-owner routes unchanged', () => {
    expect(resolveCompletedSetupRedirectPath(false, true, true)).toBe('/login')
    expect(resolveCompletedSetupRedirectPath(true, false, true)).toBe('/dashboard')
    expect(resolveCompletedSetupRedirectPath(true, true)).toBe('/admin/dashboard')
  })
  it('takes only an authenticated owner administrator to the unified gateway', () => {
    expect(resolveCompletedSetupRedirectPath(true, true, true)).toBe('/admin/everplain-gateway')
  })
})
