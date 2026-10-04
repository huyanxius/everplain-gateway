# Invoke from repository root: make -f deploy/everplain/Makefile build
.PHONY: build test-boundary
build:
	pnpm --dir frontend run build
	cd backend && GOMAXPROCS=1 GOGC=20 GOMEMLIMIT=768MiB go build -p 1 -tags embed -trimpath -ldflags="-s -w" -o bin/everplain-gateway ./cmd/server

test-boundary:
	cd backend && go test ./internal/everplain ./internal/config
