package config

import (
	"path/filepath"
	"testing"
)

func TestEverplainConfigSafety(t *testing.T) {
	valid := EverplainConfig{Enabled: true, RequestsPerMinute: 60, RequestTimeoutSeconds: 120}
	if err := valid.Validate(RunModeSimple); err != nil {
		t.Fatal(err)
	}
	if err := valid.Validate(RunModeStandard); err == nil {
		t.Fatal("standard billing must be rejected")
	}
	valid.RequestsPerMinute = 0
	if err := valid.Validate(RunModeSimple); err == nil {
		t.Fatal("unbounded requests must be rejected")
	}
	valid.Enabled = false
	if err := valid.Validate(RunModeStandard); err != nil {
		t.Fatal(err)
	}
}

func TestEverplainLightweightProfileLoads(t *testing.T) {
	resetViperWithJWTSecret(t)
	path, err := filepath.Abs("../../../deploy/everplain/config.example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_FILE", path)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Everplain.Enabled || cfg.Everplain.UpstreamEnabled || cfg.RunMode != RunModeSimple {
		t.Fatalf("unsafe gateway defaults: %+v", cfg.Everplain)
	}
	if cfg.Database.MaxOpenConns != 10 || cfg.Redis.PoolSize != 16 || cfg.Gateway.UsageRecord.WorkerCount != 4 {
		t.Fatal("lightweight pool values were ignored")
	}
	if !cfg.Security.URLAllowlist.Enabled || cfg.Security.URLAllowlist.AllowPrivateHosts || cfg.Security.URLAllowlist.AllowInsecureHTTP {
		t.Fatal("unsafe upstream URL policy")
	}
	if cfg.Pricing.RemoteURL != "" {
		t.Fatal("unexpected background remote pricing updates")
	}
}
