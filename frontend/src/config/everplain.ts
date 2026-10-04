/** This standalone distribution deliberately exposes only the gateway surface.
 * Server-side authorization and feature closure remain the security boundary.
 */
export const EVERPLAIN_GATEWAY_PROFILE = true
export const EVERPLAIN_SITE_NAME = 'Everplain Gateway'

const disabledPrefixes = [
  '/purchase', '/orders', '/payment', '/redeem', '/subscriptions', '/affiliate',
  '/admin/orders', '/admin/redeem', '/admin/subscriptions', '/admin/promo-codes', '/admin/affiliates',
  '/model-plaza', '/register', '/email-verify', '/auth', '/batch-image', '/docs/batch-image',
]
export function isGatewayRouteBlocked(path: string): boolean {
  return EVERPLAIN_GATEWAY_PROFILE && disabledPrefixes.some(prefix => path === prefix || path.startsWith(`${prefix}/`))
}

export const GATEWAY_DISABLED_FEATURES: ReadonlySet<string> = new Set([
  'subscription_enabled', 'payment_enabled', 'affiliate_enabled', 'model_plaza_enabled',
])

/** Apply before first paint and to later fetches, including stale upstream settings. */
export function applyGatewayPublicSettings<T extends object>(settings: T): T {
  return {
    ...settings,
    registration_enabled: false,
    promo_code_enabled: false,
    invitation_code_enabled: false,
    payment_enabled: false,
    subscription_enabled: false,
    affiliate_enabled: false,
    model_plaza_enabled: false,
    linuxdo_oauth_enabled: false,
    dingtalk_oauth_enabled: false,
    oidc_oauth_enabled: false,
    github_oauth_enabled: false,
    google_oauth_enabled: false,
    wechat_oauth_enabled: false,
  }
}
