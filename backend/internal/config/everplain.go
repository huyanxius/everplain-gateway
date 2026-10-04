package config

import "fmt"

// EverplainConfig narrows this fork to a server-to-server token supply gateway.
// UpstreamEnabled is an operator acknowledgement, not a credential health check.
type EverplainConfig struct {
	Enabled               bool `mapstructure:"enabled"`
	UpstreamEnabled       bool `mapstructure:"upstream_enabled"`
	RequestsPerMinute     int  `mapstructure:"requests_per_minute"`
	RequestTimeoutSeconds int  `mapstructure:"request_timeout_seconds"`
}

func (c EverplainConfig) Validate(runMode string) error {
	if !c.Enabled {
		return nil
	}
	if runMode != RunModeSimple {
		return fmt.Errorf("everplain.enabled requires run_mode=simple: Everplain owns the end-user ledger")
	}
	if c.RequestsPerMinute < 1 || c.RequestsPerMinute > 10000 {
		return fmt.Errorf("everplain.requests_per_minute must be between 1 and 10000")
	}
	if c.RequestTimeoutSeconds < 1 || c.RequestTimeoutSeconds > 600 {
		return fmt.Errorf("everplain.request_timeout_seconds must be between 1 and 600")
	}
	return nil
}
