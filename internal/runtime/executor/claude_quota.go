package executor

import (
	"net/http"
	"strconv"
	"time"
)

// ClaudeCodeQuotaInfo captures Anthropic rate limit header information.
type ClaudeCodeQuotaInfo struct {
	UnifiedStatus       string    `json:"unified_status"`
	FiveHourStatus      string    `json:"five_hour_status,omitempty"`
	FiveHourReset       int64     `json:"five_hour_reset,omitempty"`
	FiveHourUtilization float64   `json:"five_hour_utilization,omitempty"`
	SevenDayStatus      string    `json:"seven_day_status,omitempty"`
	SevenDayReset       int64     `json:"seven_day_reset,omitempty"`
	SevenDayUtilization float64   `json:"seven_day_utilization,omitempty"`
	OverageStatus       string    `json:"overage_status,omitempty"`
	OverageReset        int64     `json:"overage_reset,omitempty"`
	OverageUtilization  float64   `json:"overage_utilization,omitempty"`
	RepresentativeClaim string    `json:"representative_claim,omitempty"`
	FallbackPercentage  float64   `json:"fallback_percentage,omitempty"`
	FallbackAvailable   string    `json:"fallback_available,omitempty"`
	UnifiedReset        int64     `json:"unified_reset,omitempty"`
	LastUpdated         time.Time `json:"last_updated"`
}

// parseClaudeCodeQuotaHeaders extracts rate limit quota information from Anthropic API response headers.
// Returns nil if no quota-related headers are present.
func parseClaudeCodeQuotaHeaders(headers http.Header) *ClaudeCodeQuotaInfo {
	if headers == nil {
		return nil
	}

	unifiedStatus := headers.Get("anthropic-ratelimit-unified-status")
	if unifiedStatus == "" {
		hasAny := false
		for key := range headers {
			if len(key) > 24 && key[:24] == "Anthropic-Ratelimit-Uni" {
				hasAny = true
				break
			}
		}
		if !hasAny {
			return nil
		}
	}

	info := &ClaudeCodeQuotaInfo{
		UnifiedStatus:       unifiedStatus,
		FiveHourStatus:      headers.Get("anthropic-ratelimit-unified-5h-status"),
		FiveHourReset:       parseUnixTimestamp(headers.Get("anthropic-ratelimit-unified-5h-reset")),
		FiveHourUtilization: parseFloat(headers.Get("anthropic-ratelimit-unified-5h-utilization")),
		SevenDayStatus:      headers.Get("anthropic-ratelimit-unified-7d-status"),
		SevenDayReset:       parseUnixTimestamp(headers.Get("anthropic-ratelimit-unified-7d-reset")),
		SevenDayUtilization: parseFloat(headers.Get("anthropic-ratelimit-unified-7d-utilization")),
		OverageStatus:       headers.Get("anthropic-ratelimit-unified-overage-status"),
		OverageReset:        parseUnixTimestamp(headers.Get("anthropic-ratelimit-unified-overage-reset")),
		OverageUtilization:  parseFloat(headers.Get("anthropic-ratelimit-unified-overage-utilization")),
		RepresentativeClaim: headers.Get("anthropic-ratelimit-unified-representative-claim"),
		FallbackPercentage:  parseFloat(headers.Get("anthropic-ratelimit-unified-fallback-percentage")),
		FallbackAvailable:   headers.Get("anthropic-ratelimit-unified-fallback"),
		UnifiedReset:        parseUnixTimestamp(headers.Get("anthropic-ratelimit-unified-reset")),
		LastUpdated:         time.Now().UTC(),
	}

	return info
}

// parseUnixTimestamp converts a string Unix timestamp to int64.
// Returns 0 if the string is empty or invalid.
func parseUnixTimestamp(s string) int64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// parseFloat converts a string to float64.
// Returns 0 if the string is empty or invalid.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
