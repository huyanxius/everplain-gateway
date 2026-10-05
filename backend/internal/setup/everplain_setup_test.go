package setup

import (
	"strings"
	"testing"
)

func TestBuildPostgresDSNOmitsEmptyPassword(t *testing.T) {
	cfg := &DatabaseConfig{Host: "127.0.0.1", Port: 15432, User: "everplain_test", DBName: "everplain_gateway_test", SSLMode: "disable"}
	dsn := buildPostgresDSN(cfg, cfg.DBName)
	if strings.Contains(dsn, "password=") || !strings.Contains(dsn, "dbname=everplain_gateway_test") {
		t.Fatalf("unsafe empty-password DSN: %s", dsn)
	}
	if !strings.Contains(buildPostgresDSN(cfg, "postgres"), "dbname=postgres") {
		t.Fatal("bootstrap database override lost")
	}
}
