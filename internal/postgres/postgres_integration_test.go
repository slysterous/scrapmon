// +build integration

package postgres_test

import (
	"github.com/slysterous/scrapmon/internal/postgres"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"github.com/slysterous/scrapmon/internal/test"
	"testing"
	"time"
)

func TestNewClientError(t *testing.T) {
	connStr := "testconn"
	_, err := postgres.NewClient(connStr, 10)
	if err == nil {
		t.Fatalf("Expected NewClient to return error because of bad datasource.")
	}
}

func TestClientCreateScrap(t *testing.T) {
	db := test.DBTearUp(t)
	defer test.DBTearDown(db, t)

	client := postgres.Client{
		DB: db,
	}

	wantedScrap := scrapmon.Scrap{
		RefCode:       "00000lHB00",
		CodeCreatedAt: time.Now(),
		FileURI:       "fileuri",
	}

	id, err := client.CreateScrap(wantedScrap)
	if err != nil {
		t.Fatalf("could not create scrap err: %v", err)
	}

	if id != 1 {
		t.Errorf("Unexpected scrap id returned")
	}
}
