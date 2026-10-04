// Package everplain contains the small, credential-free integration boundary
// for the Everplain fork. Scheduling and provider protocols remain upstream.
package everplain

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ContractVersion = "2026-10-04"

// Error is stable across the upstream's provider-specific error payloads.
// Retryable never authorizes retrying a generation with an unknown outcome.
type Error struct {
	Code      string `json:"code"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Retryable bool   `json:"retryable"`
}

func ClassifyError(status int) Error {
	e := Error{Code: "internal_error", Type: "gateway_error", Message: "The gateway could not complete this request."}
	switch status {
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		e.Code, e.Message = "invalid_request", "Check the request payload and selected model."
	case http.StatusUnauthorized:
		e.Code, e.Message = "authentication_failed", "A valid server-side service key is required."
	case http.StatusForbidden:
		e.Code, e.Message = "access_denied", "This service key cannot access the requested resource."
	case http.StatusNotFound:
		e.Code, e.Message = "not_found", "This resource is unavailable in the Everplain gateway profile."
	case http.StatusRequestEntityTooLarge:
		e.Code, e.Message = "request_too_large", "The request exceeds the gateway body limit."
	case http.StatusTooManyRequests:
		e.Code, e.Message, e.Retryable = "rate_limited", "Capacity is temporarily limited; respect Retry-After.", true
	case http.StatusServiceUnavailable:
		e.Code, e.Message, e.Retryable = "provider_unavailable", "No usable upstream capacity is currently available.", true
	case http.StatusBadGateway:
		e.Code, e.Message = "upstream_error", "The upstream response was unsuccessful; its outcome may be unknown."
	case http.StatusGatewayTimeout, http.StatusRequestTimeout:
		e.Code, e.Message = "deadline_exceeded", "The request deadline elapsed; its outcome may be unknown."
	case 499:
		e.Code, e.Message = "request_cancelled", "The client cancelled the request; its outcome may be unknown."
	}
	return e
}

func Abort(c *gin.Context, status int, code, message string, retryable bool) {
	e := Error{Code: code, Type: "gateway_error", Message: message, Retryable: retryable, RequestID: c.Writer.Header().Get("X-Client-Request-ID")}
	c.Set("everplain_error", e)
	c.AbortWithStatusJSON(status, gin.H{"error": e})
}

// UsageEvent contains operational telemetry only, never end-user credits.
// Numeric row IDs are represented as strings to avoid JavaScript precision loss.
type UsageEvent struct {
	UsageID             string    `json:"usage_id"`
	RequestID           string    `json:"request_id"`
	ClientRequestID     string    `json:"client_request_id,omitempty"`
	SessionID           *string   `json:"session_id,omitempty"`
	Model               string    `json:"model"`
	UpstreamModel       *string   `json:"upstream_model,omitempty"`
	InputTokens         int       `json:"input_tokens"`
	OutputTokens        int       `json:"output_tokens"`
	CacheCreationTokens int       `json:"cache_creation_tokens"`
	CacheReadTokens     int       `json:"cache_read_tokens"`
	CostUSD             float64   `json:"cost_usd"`
	ActualCostUSD       float64   `json:"actual_cost_usd"`
	Stream              bool      `json:"stream"`
	CreatedAt           time.Time `json:"created_at"`
}

func Status(enabled, upstreamEnabled bool, rpm, timeout int) gin.H {
	state := "unconfigured"
	if upstreamEnabled {
		state = "configured_unverified"
	}
	return gin.H{
		"contract_version": ContractVersion, "enabled": enabled,
		"provider":         gin.H{"state": state, "verification": "not_performed", "agy_adapter": "antigravity_reserved"},
		"upstream_enabled": upstreamEnabled, "ledger_owner": "everplain", "budget_enforcement": "soft",
		"requests_per_minute": rpm, "request_timeout_seconds": timeout,
		"usage_consistency": "eventual", "automatic_client_retries": false,
	}
}
