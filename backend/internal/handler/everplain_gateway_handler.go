package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/everplain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// EverplainUsage exposes only the authenticated service key's telemetry, never
// account credentials, prompts, end-user identities or credit balances.
func (h *GatewayHandler) EverplainUsage(c *gin.Context) {
	key, ok := middleware.GetAPIKeyFromContext(c)
	if !ok {
		everplain.Abort(c, 401, "authentication_failed", "A valid service key is required.", false)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 || page > 100000 {
		everplain.Abort(c, 400, "invalid_request", "page must be between 1 and 100000.", false)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	if err != nil || size < 1 || size > 100 {
		everplain.Abort(c, 400, "invalid_request", "page_size must be between 1 and 100.", false)
		return
	}
	if h.usageService == nil {
		everplain.Abort(c, 503, "usage_unavailable", "Usage storage is unavailable.", true)
		return
	}
	logs, paging, err := h.usageService.ListByAPIKey(c.Request.Context(), key.ID, pagination.PaginationParams{Page: page, PageSize: size, SortOrder: pagination.SortOrderDesc})
	if err != nil {
		everplain.Abort(c, 503, "usage_unavailable", "Usage storage is unavailable.", true)
		return
	}
	events := make([]everplain.UsageEvent, 0, len(logs))
	for _, entry := range logs {
		events = append(events, everplainUsageEvent(entry))
	}
	total := int64(0)
	if paging != nil {
		total = paging.Total
	}
	c.JSON(http.StatusOK, gin.H{"contract_version": everplain.ContractVersion, "ledger_owner": "everplain", "consistency": "eventual", "items": events, "page": page, "page_size": size, "total": total})
}
func everplainUsageEvent(entry service.UsageLog) everplain.UsageEvent {
	model := entry.RequestedModel
	if model == "" {
		model = entry.Model
	}
	clientRequestID := ""
	if strings.HasPrefix(entry.RequestID, "client:") {
		clientRequestID = strings.TrimPrefix(entry.RequestID, "client:")
	}
	return everplain.UsageEvent{ClientRequestID: clientRequestID, UsageID: strconv.FormatInt(entry.ID, 10), RequestID: entry.RequestID, SessionID: entry.SessionID, Model: model, UpstreamModel: entry.UpstreamModel, InputTokens: entry.InputTokens, OutputTokens: entry.OutputTokens, CacheCreationTokens: entry.CacheCreationTokens, CacheReadTokens: entry.CacheReadTokens, CostUSD: entry.TotalCost, ActualCostUSD: entry.ActualCost, Stream: entry.Stream, CreatedAt: entry.CreatedAt}
}
