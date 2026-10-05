package admin

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	everplainGatewayOwnerEmail       = "huyanxius@gmail.com"
	everplainGatewayDefaultPageSize  = 25
	everplainGatewayMaxPageSize      = 100
	everplainGatewayMaxPage          = 100000
	everplainGatewayProviderCostNote = "Gateway actual_cost/account_cost is internal billing and allocation data, not provider-invoiced cost."
)

// This projection deliberately does not embed Account or an existing account DTO.
// Adding fields to those types must never expand the inventory's disclosure.
type everplainGatewayUpstream struct {
	ID           int64    `json:"id"`
	Alias        string   `json:"alias"`
	Platform     string   `json:"platform"`
	Type         string   `json:"type"`
	Status       string   `json:"status"`
	Schedulable  bool     `json:"schedulable"`
	GroupIDs     []int64  `json:"group_ids"`
	MappedModels []string `json:"mapped_models"`
}

type everplainGatewayProviderCharge struct {
	Status   string   `json:"status"`
	Amount   *float64 `json:"amount"`
	Currency *string  `json:"currency"`
	Source   *string  `json:"source"`
	Reason   string   `json:"reason"`
}

type everplainGatewayUpstreamsResponse struct {
	Items          []everplainGatewayUpstream     `json:"items"`
	Count          int                            `json:"count"`
	Total          int64                          `json:"total"`
	Page           int                            `json:"page"`
	PageSize       int                            `json:"page_size"`
	ProviderCharge everplainGatewayProviderCharge `json:"provider_charge"`
}

// ListEverplainGatewayUpstreams is a local, read-only inventory. It deliberately
// bypasses account usage/billing resolvers and never probes upstream providers.
// GET /api/v1/admin/everplain-gateway/upstreams
func (h *AccountHandler) ListEverplainGatewayUpstreams(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	subject, authenticated := middleware.GetAuthSubjectFromContext(c)
	role, hasRole := middleware.GetUserRoleFromContext(c)
	if c.GetString("auth_method") != service.AuditAuthMethodJWT ||
		!authenticated || subject.UserID <= 0 || !hasRole || role != domain.RoleAdmin ||
		!strings.EqualFold(c.GetString(middleware.ContextKeyAuthEmail), everplainGatewayOwnerEmail) {
		response.ErrorWithDetails(c, http.StatusForbidden, "Owner admin session required", "EVERPLAIN_GATEWAY_OWNER_REQUIRED", nil)
		return
	}

	page, validPage := everplainGatewayPageParameter(c, "page", 1, everplainGatewayMaxPage)
	pageSize, validPageSize := everplainGatewayPageParameter(c, "page_size", everplainGatewayDefaultPageSize, everplainGatewayMaxPageSize)
	if !validPage || !validPageSize {
		response.ErrorWithDetails(c, http.StatusBadRequest, "Invalid inventory pagination", "EVERPLAIN_GATEWAY_INVALID_PAGINATION", nil)
		return
	}
	if h.adminService == nil {
		response.InternalError(c, "Upstream inventory unavailable")
		return
	}

	accounts, total, err := h.adminService.ListAccounts(c.Request.Context(), page, pageSize, "", "", "", "", 0, "", "id", "asc")
	if err != nil {
		// Repository errors can contain private account or connection details.
		response.InternalError(c, "Upstream inventory unavailable")
		return
	}

	items := make([]everplainGatewayUpstream, 0, len(accounts))
	for i := range accounts {
		account := &accounts[i]
		groupIDs := append([]int64{}, account.GroupIDs...)
		items = append(items, everplainGatewayUpstream{
			ID: account.ID, Alias: "upstream-" + strconv.FormatInt(account.ID, 10),
			Platform: account.Platform, Type: account.Type, Status: account.Status,
			Schedulable: account.Schedulable, GroupIDs: groupIDs,
			MappedModels: everplainGatewayMappedModels(account.Credentials["model_mapping"]),
		})
	}
	response.Success(c, everplainGatewayUpstreamsResponse{
		Items: items, Count: len(items), Total: total, Page: page, PageSize: pageSize,
		ProviderCharge: everplainGatewayProviderCharge{Status: "unavailable", Reason: everplainGatewayProviderCostNote},
	})
}

func everplainGatewayPageParameter(c *gin.Context, name string, defaultValue, maximum int) (int, bool) {
	values, present := c.Request.URL.Query()[name]
	if !present {
		return defaultValue, true
	}
	if len(values) != 1 {
		return 0, false
	}
	value, err := strconv.Atoi(values[0])
	return value, err == nil && value >= 1 && value <= maximum
}

// Only existing mapping values are exposed. There are no mapping keys, inferred
// defaults, base URLs, provider requests, or credential-presence diagnostics.
func everplainGatewayMappedModels(raw any) []string {
	models := make(map[string]struct{})
	add := func(value string) {
		if everplainGatewaySafeModelName(value) {
			models[value] = struct{}{}
		}
	}
	switch mapping := raw.(type) {
	case map[string]any:
		for _, value := range mapping {
			if name, ok := value.(string); ok {
				add(name)
			}
		}
	case map[string]string:
		for _, value := range mapping {
			add(value)
		}
	}
	result := make([]string, 0, len(models))
	for name := range models {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func everplainGatewaySafeModelName(value string) bool {
	if value == "" || len(value) > 200 || strings.Contains(value, "://") {
		return false
	}
	lower := strings.ToLower(value)
	for _, prefix := range []string{"sk-", "sk_", "rk-", "rk_", "pk-", "pk_", "sess-", "ghp_", "gho_", "github_pat_", "akia", "asia", "aiza", "ya29.", "hf_", "gsk_", "eyj"} {
		if strings.HasPrefix(lower, prefix) {
			return false
		}
	}
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			continue
		}
		if i > 0 && strings.ContainsRune("._:/-", rune(ch)) {
			continue
		}
		return false
	}
	return true
}
