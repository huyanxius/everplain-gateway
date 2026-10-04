package everplain

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func hasPathPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

// DisabledRoute also covers direct HTTP calls; hiding a navigation item alone
// must not reactivate the upstream's sales or affiliate functionality.
func DisabledRoute(path string) bool {
	for _, p := range []string{
		"/api/v1/payment", "/api/v1/admin/payment", "/api/v1/redeem", "/api/v1/subscriptions",
		"/api/v1/admin/redeem-codes", "/api/v1/admin/promo-codes", "/api/v1/admin/subscriptions",
		"/api/v1/admin/affiliates", "/api/v1/admin/affiliate", "/api/v1/user/aff",
		"/api/v1/model-plaza", "/api/v1/auth/register", "/api/v1/auth/validate-promo-code",
		"/api/v1/auth/validate-invitation-code", "/api/v1/admin/system/update", "/api/v1/admin/system/rollback", "/api/v1/auth/oauth", "/v1/sub2api/billing",
	} {
		if hasPathPrefix(path, p) {
			return true
		}
	}
	if strings.HasPrefix(path, "/api/v1/admin/users/") && (strings.HasSuffix(path, "/balance") || strings.HasSuffix(path, "/balance-history") || strings.HasSuffix(path, "/subscriptions")) {
		return true
	}
	if strings.HasPrefix(path, "/api/v1/admin/openai/accounts/") && strings.Contains(path, "/referrals") {
		return true
	}
	return strings.HasPrefix(path, "/api/v1/admin/groups/") && strings.HasSuffix(path, "/subscriptions")
}

// IsGatewayPath recognizes all original provider entry points so aliases cannot
// bypass the deliberately narrow Everplain text-only contract.
func IsGatewayPath(path string) bool {
	for _, prefix := range []string{"/backend-api", "/api/v3", "/v3", "/contents", "/alpha", "/embeddings", "/images", "/videos", "/tts", "/stt", "/custom-voices", "/realtime", "/web_search", "/x_search", "/v1", "/v1beta", "/antigravity", "/responses", "/chat/completions", "/messages", "/models", "/gemini", "/openai", "/claude", "/api/v1beta"} {
		if hasPathPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func AllowedGatewayRoute(method, path string) bool {
	if method == http.MethodGet {
		return path == "/v1/models" || strings.HasPrefix(path, "/v1/models/") || path == "/v1/usage" || path == "/v1/everplain/usage"
	}
	return method == http.MethodPost && (path == "/v1/messages" || path == "/v1/responses" || path == "/v1/chat/completions")
}

func ProfileGuard(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if enabled && (DisabledRoute(c.Request.URL.Path) || (IsGatewayPath(c.Request.URL.Path) && !AllowedGatewayRoute(c.Request.Method, c.Request.URL.Path))) {
			Abort(c, http.StatusNotFound, "feature_disabled", "This endpoint is disabled in the Everplain gateway profile.", false)
			return
		}
		c.Next()
	}
}

// GenerationGuard runs after API-key authentication. It does not fabricate
// readiness or configure an account. No provider call is possible by default.
func GenerationGuard(upstreamEnabled bool, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}
		if !upstreamEnabled {
			Abort(c, http.StatusServiceUnavailable, "provider_not_configured", "No upstream has been enabled. Configure and verify an authorized provider first.", false)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		if ctx.Err() != nil {
			Abort(c, 499, "request_cancelled", "The client cancelled the request.", false)
			return
		}
		c.Next()
		if !c.Writer.Written() && ctx.Err() != nil {
			status := 499
			if ctx.Err() == context.DeadlineExceeded {
				status = http.StatusGatewayTimeout
			}
			e := ClassifyError(status)
			Abort(c, status, e.Code, e.Message, false)
		}
	}
}
