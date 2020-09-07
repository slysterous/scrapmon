package postgres_test

import (
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"github.com/slysterous/print-scrape/internal/postgres"
	"github.com/slysterous/print-scrape/internal/test"
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
	db := test.TearUp(t)
	defer test.TearDown(db, t)

	client := postgres.Client{
		DB: db,
	}

	wantedScrap := printscrape.ScreenShot{
		RefCode:       "00000lHB00",
		CodeCreatedAt: time.Now(),
		FileURI:       "fileuri",
	}

	id, err := client.CreateScreenShot(wantedScrap)
	if err != nil {
		t.Fatalf("could not create scrap err: %v", err)
	}

	if id != 1 {
		t.Errorf("Unexpected scrap id returned")
	}
}
