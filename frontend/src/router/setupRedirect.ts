export function resolveCompletedSetupRedirectPath(isAuthenticated: boolean, isAdmin: boolean, gatewayOwner = false): string {
  if (!isAuthenticated) {
    return '/login'
  }

  return isAdmin ? (gatewayOwner ? '/admin/everplain-gateway' : '/admin/dashboard') : '/dashboard'
}
