package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/everplain"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Enumerate the actual upstream route table, not a hand-maintained alias list.
// Every provider entry point outside the approved contract must be fail-closed.
func TestEverplainAllRegisteredGatewayRoutesAreGuarded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(everplain.NormalizeErrors(true), everplain.ProfileGuard(true))
	cfg := &config.Config{Everplain: config.EverplainConfig{Enabled: true}, Gateway: config.GatewayConfig{MaxBodySize: 1024 * 1024, TextMaxBodySize: 1024 * 1024}}
	authCalls := 0
	auth := middleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		authCalls++
		c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{ID: 1})
		c.Next()
	})
	RegisterGatewayRoutes(r, &handler.Handlers{Gateway: &handler.GatewayHandler{}, OpenAIGateway: &handler.OpenAIGatewayHandler{}, AsyncImage: handler.NewAsyncImageHandler(nil, nil)}, auth, nil, nil, nil, nil, nil, cfg, everplain.GenerationGuard(false, time.Second))
	for _, route := range r.Routes() {
		if everplain.AllowedGatewayRoute(route.Method, route.Path) {
			continue
		}
		path := route.Path
		for _, part := range strings.Split(path, "/") {
			if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
				path = strings.Replace(path, part, "test", 1)
			}
		}
		before := authCalls
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(route.Method, path, nil))
		if w.Code != http.StatusNotFound || before != authCalls {
			t.Errorf("unguarded upstream alias %s %s status=%d auth_calls=%d", route.Method, path, w.Code, authCalls-before)
		}
	}
	for _, path := range []string{"/v1/messages", "/v1/responses", "/v1/chat/completions"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, nil))
		if w.Code != 503 || !strings.Contains(w.Body.String(), "provider_not_configured") {
			t.Errorf("unconfigured %s: %d %s", path, w.Code, w.Body.String())
		}
	}
}
