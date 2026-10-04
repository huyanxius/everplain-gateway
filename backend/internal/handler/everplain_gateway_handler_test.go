package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEverplainUsageEventSafeProjection(t *testing.T) {
	model := "mapped-model"
	session := "opaque-operation-7"
	event := everplainUsageEvent(service.UsageLog{ID: 9007199254740993, RequestID: "client:gateway-correlation", SessionID: &session, RequestedModel: "public-model", Model: "legacy-model", UpstreamModel: &model, InputTokens: 12, OutputTokens: 3, CacheReadTokens: 2, UserID: 91, AccountID: 92, APIKeyID: 93, TotalCost: 0.001, ActualCost: 0.002})
	b, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if event.ClientRequestID != "gateway-correlation" || event.UsageID != "9007199254740993" || event.Model != "public-model" || event.InputTokens != 12 {
		t.Fatalf("wrong event %+v", event)
	}
	for _, name := range []string{"user_id", "account_id", "api_key_id", "balance", "credits"} {
		if strings.Contains(string(b), name) {
			t.Fatalf("unsafe field %s", name)
		}
	}
	fallback := everplainUsageEvent(service.UsageLog{Model: "historical"})
	if fallback.Model != "historical" {
		t.Fatal("legacy model missing")
	}
}

type everplainUsageRepoStub struct {
	service.UsageLogRepository
	gotKey int64
	fail   bool
}

func (s *everplainUsageRepoStub) ListByAPIKey(_ context.Context, key int64, p pagination.PaginationParams) ([]service.UsageLog, *pagination.PaginationResult, error) {
	s.gotKey = key
	if s.fail {
		return nil, nil, errors.New("private database error")
	}
	return []service.UsageLog{{ID: 42, APIKeyID: key, Model: "test"}}, &pagination.PaginationResult{Total: 1, Page: p.Page, PageSize: p.PageSize}, nil
}
func TestEverplainUsageRequiresKeyAndScopesRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &everplainUsageRepoStub{}
	h := &GatewayHandler{usageService: service.NewUsageService(repo, nil, nil, nil)}
	r := gin.New()
	r.GET("/v1/everplain/usage", func(c *gin.Context) {
		if c.GetHeader("X-Test-Key") == "yes" {
			c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{ID: 17})
		}
	}, h.EverplainUsage)
	call := func(path string, auth bool) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		if auth {
			req.Header.Set("X-Test-Key", "yes")
		}
		r.ServeHTTP(w, req)
		return w
	}
	if w := call("/v1/everplain/usage", false); w.Code != 401 || repo.gotKey != 0 {
		t.Fatalf("missing key accessed records: %d", w.Code)
	}
	if w := call("/v1/everplain/usage", true); w.Code != 200 || repo.gotKey != 17 || !strings.Contains(w.Body.String(), `"usage_id":"42"`) {
		t.Fatalf("wrong scope: %d %s", w.Code, w.Body.String())
	}
	if w := call("/v1/everplain/usage?page_size=1000", true); w.Code != 400 {
		t.Fatal("unbounded pagination accepted")
	}
	repo.fail = true
	if w := call("/v1/everplain/usage", true); w.Code != 503 || strings.Contains(w.Body.String(), "private database error") {
		t.Fatal("storage failure hidden or leaked")
	}
}
