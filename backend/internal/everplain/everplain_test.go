package everplain

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func engine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Header("X-Client-Request-ID", "test-correlation"); c.Next() })
	return r
}
func request(r http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}
func errorBody(t *testing.T, w *httptest.ResponseRecorder) Error {
	t.Helper()
	var body struct {
		Error Error `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.Error
}
func TestProfileBlocksSalesAndAliases(t *testing.T) {
	for _, path := range []string{"/api/v1/payment/webhook/stripe", "/api/v1/admin/payment/config", "/api/v1/admin/promo-codes", "/api/v1/admin/affiliates", "/api/v1/user/aff/transfer", "/api/v1/auth/register", "/api/v1/auth/oauth/google/callback", "/api/v1/admin/users/42/balance", "/api/v1/admin/groups/1/subscriptions", "/v1/sub2api/billing", "/backend-api/codex/responses", "/api/v3/contents/generations/tasks", "/antigravity/v1/messages", "/responses", "/images/generations", "/v1/images/generations", "/embeddings"} {
		t.Run(path, func(t *testing.T) {
			r := engine()
			r.Use(ProfileGuard(true))
			called := false
			r.Any(path, func(c *gin.Context) { called = true; c.Status(200) })
			w := request(r, "POST", path)
			if w.Code != 404 || called {
				t.Fatalf("code=%d called=%v", w.Code, called)
			}
		})
	}
	for _, path := range []string{"/api/v1/admin/accounts", "/api/v1/admin/groups", "/api/v1/admin/usage", "/api/v1/keys", "/api/v1/auth/login", "/api/v1/admin/antigravity/auth-url"} {
		if DisabledRoute(path) {
			t.Fatalf("core route blocked: %s", path)
		}
	}
}
func TestUnconfiguredNeverCallsUpstream(t *testing.T) {
	r := engine()
	r.Use(NormalizeErrors(true))
	r.POST("/v1/responses", GenerationGuard(false, time.Second), func(c *gin.Context) { t.Fatal("unexpected provider call") })
	w := request(r, "POST", "/v1/responses")
	e := errorBody(t, w)
	if w.Code != 503 || e.Code != "provider_not_configured" || e.Retryable || e.RequestID != "test-correlation" {
		t.Fatalf("%d %+v", w.Code, e)
	}
}
func TestErrorNormalizationDoesNotLeakOrRetryUnknownOutcome(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 413, 429, 499, 500, 502, 503, 504} {
		r := engine()
		r.Use(NormalizeErrors(true))
		r.POST("/v1/messages", func(c *gin.Context) {
			c.Header("Content-Length", "999")
			c.JSON(status, gin.H{"error": "secret-provider-account-and-prompt"})
		})
		w := request(r, "POST", "/v1/messages")
		e := errorBody(t, w)
		if w.Code != status || strings.Contains(w.Body.String(), "secret") || e.RequestID == "" || w.Header().Get("Content-Length") != "" {
			t.Fatalf("invalid normalization %d %s", w.Code, w.Body.String())
		}
		if status != 429 && status != 503 && e.Retryable {
			t.Fatalf("unsafe retry status %d", status)
		}
	}
}
func TestSuccessfulSSEIsUnbufferedAndNeverRetried(t *testing.T) {
	r := engine()
	r.Use(NormalizeErrors(true))
	var calls atomic.Int32
	r.POST("/v1/messages", GenerationGuard(true, time.Second), func(c *gin.Context) {
		calls.Add(1)
		c.Header("Content-Type", "text/event-stream")
		_, _ = c.Writer.WriteString("data: {\"text\":\"hello\"}\n\n")
		c.Writer.Flush()
		_, _ = c.Writer.WriteString("event: error\ndata: {\"type\":\"upstream_error\"}\n\n")
	})
	w := request(r, "POST", "/v1/messages")
	if calls.Load() != 1 || !w.Flushed || !strings.Contains(w.Body.String(), "hello") || !strings.Contains(w.Body.String(), "event: error") {
		t.Fatalf("SSE changed %s", w.Body.String())
	}
}
func TestCancellationStopsBeforeProviderAndPropagates(t *testing.T) {
	r := engine()
	called := false
	r.POST("/v1/messages", GenerationGuard(true, time.Second), func(c *gin.Context) { called = true })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/v1/messages", nil).WithContext(ctx))
	if called || w.Code != 499 {
		t.Fatalf("cancelled request ran: %v %d", called, w.Code)
	}
	r = engine()
	r.POST("/v1/messages", GenerationGuard(true, 5*time.Millisecond), func(c *gin.Context) { <-c.Request.Context().Done() })
	w = request(r, "POST", "/v1/messages")
	if w.Code != 504 {
		t.Fatalf("deadline: %d", w.Code)
	}
}
func TestRateLimitIsSharedExpiresAndFailsClosed(t *testing.T) {
	mini := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mini.Addr(), MaxRetries: -1})
	t.Cleanup(func() { _ = client.Close() })
	makeRouter := func() *gin.Engine {
		r := engine()
		r.POST("/v1/messages", RateLimit(client, 2, func(*gin.Context) int64 { return 7 }), func(c *gin.Context) { c.Status(204) })
		return r
	}
	a, b := makeRouter(), makeRouter()
	if request(a, "POST", "/v1/messages").Code != 204 || request(b, "POST", "/v1/messages").Code != 204 {
		t.Fatal("first requests failed")
	}
	w := request(a, "POST", "/v1/messages")
	if w.Code != 429 || w.Header().Get("Retry-After") != "60" {
		t.Fatalf("limit: %d %v", w.Code, w.Header())
	}
	mini.FastForward(time.Minute)
	if request(b, "POST", "/v1/messages").Code != 204 {
		t.Fatal("window did not expire")
	}
	mini.Close()
	if request(a, "POST", "/v1/messages").Code != 503 {
		t.Fatal("redis failure must fail closed")
	}
}
func TestStatusDoesNotClaimCredentialReadiness(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		s := Status(true, enabled, 60, 120)
		provider := s["provider"].(gin.H)
		if provider["verification"] != "not_performed" || s["budget_enforcement"] != "soft" || s["ledger_owner"] != "everplain" {
			t.Fatalf("invalid status %+v", s)
		}
	}
}

func TestDeadlineCancelsUpstreamTransportWithoutRetry(t *testing.T) {
	var calls atomic.Int32
	cancelled := make(chan struct{}, 1)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		<-r.Context().Done()
		cancelled <- struct{}{}
	}))
	defer upstream.Close()
	r := engine()
	r.Use(NormalizeErrors(true))
	r.POST("/v1/messages", GenerationGuard(true, 20*time.Millisecond), func(c *gin.Context) {
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, upstream.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := upstream.Client().Do(req)
		if err == nil {
			_ = resp.Body.Close()
			t.Error("expected context cancellation")
		}
	})
	w := request(r, "POST", "/v1/messages")
	if w.Code != 504 {
		t.Fatalf("deadline status %d: %s", w.Code, w.Body.String())
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("upstream request context was not cancelled")
	}
	if calls.Load() != 1 {
		t.Fatalf("unexpected generation retries: %d", calls.Load())
	}
}

func TestRateLimitConcurrentAdmissionIsAtomic(t *testing.T) {
	mini := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mini.Addr()})
	defer func() { _ = client.Close() }()
	r := engine()
	r.POST("/v1/messages", RateLimit(client, 3, func(*gin.Context) int64 { return 19 }), func(c *gin.Context) { c.Status(204) })
	results := make(chan int, 20)
	for range 20 {
		go func() { results <- request(r, "POST", "/v1/messages").Code }()
	}
	accepted, limited := 0, 0
	for range 20 {
		switch code := <-results; code {
		case 204:
			accepted++
		case 429:
			limited++
		default:
			t.Fatalf("unexpected status %d", code)
		}
	}
	if accepted != 3 || limited != 17 {
		t.Fatalf("atomic admission failed: accepted=%d limited=%d", accepted, limited)
	}
}
