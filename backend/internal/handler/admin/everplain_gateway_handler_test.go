package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type everplainInventoryServiceStub struct {
	service.AdminService
	accounts []service.Account
	total    int64
	err      error
	calls    int
	page     int
	pageSize int
	filters  []string
	groupID  int64
}

func (s *everplainInventoryServiceStub) ListAccounts(_ context.Context, page, pageSize int, platform, accountType, status, search string, groupID int64, privacyMode, sortBy, sortOrder string) ([]service.Account, int64, error) {
	s.calls++
	s.page, s.pageSize = page, pageSize
	s.filters = []string{platform, accountType, status, search, privacyMode, sortBy, sortOrder}
	s.groupID = groupID
	return s.accounts, s.total, s.err
}

type everplainInventoryIdentity struct {
	method  string
	email   string
	role    string
	subject any
}

func everplainOwnerIdentity() everplainInventoryIdentity {
	return everplainInventoryIdentity{
		method: service.AuditAuthMethodJWT, email: everplainGatewayOwnerEmail,
		role: domain.RoleAdmin, subject: middleware.AuthSubject{UserID: 1},
	}
}

func everplainInventoryRequest(t *testing.T, svc *everplainInventoryServiceStub, identity everplainInventoryIdentity, query string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("auth_method", identity.method)
		c.Set(middleware.ContextKeyAuthEmail, identity.email)
		c.Set(string(middleware.ContextKeyUserRole), identity.role)
		if identity.subject != nil {
			c.Set(string(middleware.ContextKeyUser), identity.subject)
		}
	})
	h := &AccountHandler{adminService: svc}
	router.GET("/api/v1/admin/everplain-gateway/upstreams", h.ListEverplainGatewayUpstreams)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/admin/everplain-gateway/upstreams"+query, nil))
	return recorder
}

func TestEverplainGatewayInventoryRequiresOwnerJWTIdentity(t *testing.T) {
	cases := []struct {
		name   string
		change func(*everplainInventoryIdentity)
	}{
		{"other-admin", func(i *everplainInventoryIdentity) { i.email = "other@example.com" }},
		{"admin-api-key-owner", func(i *everplainInventoryIdentity) { i.method = service.AuditAuthMethodAdminAPIKey }},
		{"missing-auth-method", func(i *everplainInventoryIdentity) { i.method = "" }},
		{"missing-email", func(i *everplainInventoryIdentity) { i.email = "" }},
		{"email-whitespace", func(i *everplainInventoryIdentity) { i.email = " " + i.email }},
		{"missing-subject", func(i *everplainInventoryIdentity) { i.subject = nil }},
		{"wrong-subject-type", func(i *everplainInventoryIdentity) { i.subject = map[string]any{"UserID": 1} }},
		{"zero-user-id", func(i *everplainInventoryIdentity) { i.subject = middleware.AuthSubject{} }},
		{"negative-user-id", func(i *everplainInventoryIdentity) { i.subject = middleware.AuthSubject{UserID: -1} }},
		{"non-admin-owner", func(i *everplainInventoryIdentity) { i.role = domain.RoleUser }},
		{"missing-role", func(i *everplainInventoryIdentity) { i.role = "" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &everplainInventoryServiceStub{}
			identity := everplainOwnerIdentity()
			tc.change(&identity)
			recorder := everplainInventoryRequest(t, svc, identity, "")
			require.Equal(t, http.StatusForbidden, recorder.Code)
			require.Zero(t, svc.calls, "denied identity must not read accounts")
			require.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
		})
	}

	identity := everplainOwnerIdentity()
	identity.email = strings.ToUpper(identity.email)
	svc := &everplainInventoryServiceStub{}
	recorder := everplainInventoryRequest(t, svc, identity, "")
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, svc.calls)
}

func TestEverplainGatewayInventoryRedactsPrivateAccountFields(t *testing.T) {
	// Never include fixture secrets in assertion output, even when a test fails.
	privateValues := []string{
		"private" + "-account-name", "private" + "-account-notes", "sk-" + "fixture-secret",
		"private" + "-email@example.com", "private" + "-error-message", "private" + "-extra-value",
		"https://provider.invalid/" + "v1?api_key=fixture-secret", "private" + "-mapping-key",
	}
	notes := privateValues[1]
	svc := &everplainInventoryServiceStub{
		accounts: []service.Account{{
			ID: 41, Name: privateValues[0], Notes: &notes, ErrorMessage: privateValues[4],
			Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey, Status: service.StatusActive,
			Schedulable: true, GroupIDs: []int64{3, 8},
			Credentials: map[string]any{
				"api_key": privateValues[2], "email": privateValues[3], "base_url": privateValues[6],
				"model_mapping": map[string]any{
					privateValues[7]: "gpt-4.1", "duplicate": "gpt-4.1", "vendor-alias": "meta-llama/Llama-3.3-70B",
					"not-a-model": privateValues[2], "url-value": privateValues[6], "non-string": true,
				},
			},
			Extra: map[string]any{"private": privateValues[5], "actual_cost": 12.3, "account_cost": 4.5},
		}},
		total: 141,
	}
	recorder := everplainInventoryRequest(t, svc, everplainOwnerIdentity(), "?page=2&page_size=10")
	for _, value := range privateValues {
		if strings.Contains(recorder.Body.String(), value) {
			t.Fatal("inventory exposed a private fixture value")
		}
	}
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
	var payload struct {
		Code int `json:"code"`
		Data struct {
			Items          []map[string]any `json:"items"`
			Count          int              `json:"count"`
			Total          int64            `json:"total"`
			Page           int              `json:"page"`
			PageSize       int              `json:"page_size"`
			ProviderCharge map[string]any   `json:"provider_charge"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
	require.Zero(t, payload.Code)
	require.Equal(t, 1, payload.Data.Count)
	require.Equal(t, int64(141), payload.Data.Total)
	require.Equal(t, 2, payload.Data.Page)
	require.Equal(t, 10, payload.Data.PageSize)
	require.Len(t, payload.Data.Items, 1)
	item := payload.Data.Items[0]
	keys := make([]string, 0, len(item))
	for key := range item {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	require.Equal(t, []string{"alias", "group_ids", "id", "mapped_models", "platform", "schedulable", "status", "type"}, keys)
	require.Equal(t, "upstream-41", item["alias"])
	require.Equal(t, []any{float64(3), float64(8)}, item["group_ids"])
	require.Equal(t, []any{"gpt-4.1", "meta-llama/Llama-3.3-70B"}, item["mapped_models"])
	charge := payload.Data.ProviderCharge
	require.Equal(t, "unavailable", charge["status"])
	for _, field := range []string{"amount", "currency", "source"} {
		require.Contains(t, charge, field)
		require.Nil(t, charge[field])
	}
	require.Equal(t, everplainGatewayProviderCostNote, charge["reason"])
	require.Equal(t, []string{"", "", "", "", "", "id", "asc"}, svc.filters)
	require.Zero(t, svc.groupID)
	require.Equal(t, 2, svc.page)
	require.Equal(t, 10, svc.pageSize)
}

func TestEverplainGatewayInventoryPaginationIsBounded(t *testing.T) {
	for _, query := range []string{
		"?page=0", "?page=-1", "?page=100001", "?page=nope", "?page=", "?page=1&page=2",
		"?page_size=0", "?page_size=-1", "?page_size=101", "?page_size=nope", "?page_size=", "?page_size=1&page_size=2",
		"?page=99999999999999999999999999",
	} {
		t.Run(query, func(t *testing.T) {
			svc := &everplainInventoryServiceStub{}
			recorder := everplainInventoryRequest(t, svc, everplainOwnerIdentity(), query)
			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, svc.calls)
		})
	}
	for _, tc := range []struct {
		query    string
		page     int
		pageSize int
	}{
		{"", 1, 25}, {"?page=100000&page_size=100", 100000, 100},
	} {
		svc := &everplainInventoryServiceStub{}
		recorder := everplainInventoryRequest(t, svc, everplainOwnerIdentity(), tc.query)
		require.Equal(t, http.StatusOK, recorder.Code)
		require.Equal(t, tc.page, svc.page)
		require.Equal(t, tc.pageSize, svc.pageSize)
		require.Contains(t, recorder.Body.String(), `"items":[]`)
	}
}

func TestEverplainGatewayInventoryDoesNotExposeRepositoryErrors(t *testing.T) {
	privateDetail := "private" + "-repository-detail"
	svc := &everplainInventoryServiceStub{err: errors.New(privateDetail)}
	recorder := everplainInventoryRequest(t, svc, everplainOwnerIdentity(), "")
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	if strings.Contains(recorder.Body.String(), privateDetail) {
		t.Fatal("inventory exposed private repository error details")
	}
	require.Equal(t, 1, svc.calls)
}

func TestEverplainGatewayMappedModelsAcceptsOnlySafeValues(t *testing.T) {
	require.Equal(t, []string{"claude-sonnet-4-5", "gpt-4.1"}, everplainGatewayMappedModels(map[string]string{
		"source-1": "gpt-4.1", "source-2": "claude-sonnet-4-5", "duplicate": "gpt-4.1",
	}))
	require.Equal(t, []string{}, everplainGatewayMappedModels(nil))
	for _, value := range []string{
		"", " model ", "https://provider.invalid/v1", "model?key=value", "model&key=value", "model@provider.invalid",
		"model%3Fkey=value", "model\nvalue", "/model", strings.Repeat("m", 201),
		"sk-" + "test-value", "sk_" + "test-value", "ghp_" + "test-value", "eyJ" + "test-value",
		"AIza" + "test-value", "ya29." + "test-value", "hf_" + "test-value", "gsk_" + "test-value",
	} {
		if everplainGatewaySafeModelName(value) {
			t.Fatal("unsafe model identifier was accepted")
		}
	}
}
