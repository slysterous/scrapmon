#!/bin/sh
set -e

./scripts/gofmtcheck.sh

# lint
echo "==> Checking lint"
#golint -set_exit_status=1 `go list -mod=vendor ./...`

# vet
echo "==> Checking vet"
go vet -mod vendor ./...

echo "==> Verifying modules"
go mod verify

# test
echo "==> Running all tests with race detector and coverage (this will take a while)"
go test ./... -race -p 1 -cover -v -tags=integration -mod vendor