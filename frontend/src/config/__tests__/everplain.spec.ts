import { describe, expect, it } from 'vitest'
import { isGatewayRouteBlocked, GATEWAY_DISABLED_FEATURES, applyGatewayPublicSettings } from '../everplain'

describe('standalone gateway surface', () => {
  it.each(['/purchase', '/orders', '/payment/result', '/payment/stripe-popup', '/redeem', '/subscriptions', '/affiliate', '/admin/orders/plans', '/admin/redeem', '/admin/subscriptions', '/admin/promo-codes', '/admin/affiliates/invites', '/model-plaza', '/register', '/email-verify', '/auth/callback', '/auth/wechat/payment/callback', '/batch-image'])('closes %s before its view loads', path => {
    expect(isGatewayRouteBlocked(path)).toBe(true)
  })
  it.each(['/setup', '/login', '/admin/dashboard', '/admin/accounts', '/admin/groups', '/keys', '/usage', '/admin/usage', '/admin/settings', '/profile', '/legal/terms'])('keeps %s reachable under upstream auth', path => {
    expect(isGatewayRouteBlocked(path)).toBe(false)
  })
  it('closes stale public settings without changing unrelated settings', () => {
    const input = { registration_enabled: true, payment_enabled: true, google_oauth_enabled: true, site_name: 'Test' }
    const output = applyGatewayPublicSettings(input)
    expect(output).toMatchObject({ registration_enabled: false, payment_enabled: false, google_oauth_enabled: false, site_name: 'Test' })
    expect(input.registration_enabled).toBe(true)
  })
  it('does not overmatch prefix lookalikes', () => {
    expect(isGatewayRouteBlocked('/authorization')).toBe(false)
    expect(isGatewayRouteBlocked('/orders-archive')).toBe(false)
  })
  it('disables background commerce feature consumers even before public settings load', () => {
    expect([...GATEWAY_DISABLED_FEATURES]).toEqual(expect.arrayContaining(['subscription_enabled', 'payment_enabled', 'affiliate_enabled', 'model_plaza_enabled']))
  })
})
