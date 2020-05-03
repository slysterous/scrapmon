package postgres_test

import (
	"github.com/slysterous/print-scrape/internal/postgres"
	"testing"
)

func TestNewClientError(t *testing.T) {
	connStr := "testconn"
	_, err := postgres.NewClient(connStr, 10)
	if err == nil {
		t.Fatalf("Expected NewClient to return error because of bad datasource.")
	}
}
