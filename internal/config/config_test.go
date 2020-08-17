package config_test

import (
	"os"
	"testing"

	"github.com/slysterous/print-scrape/internal/config"
)

func TestFromEnv(t *testing.T) {
	setEnv(t, "PRINT_SCRAPE_DB_HOST", "host")
	setEnv(t, "PRINT_SCRAPE_DB_NAME", "db")
	setEnv(t, "PRINT_SCRAPE_DB_PORT", "5000")
	setEnv(t, "PRINT_SCRAPE_DB_USER", "dbuser")
	setEnv(t, "PRINT_SCRAPE_DB_PASSWORD", "password")
	setEnv(t, "MAX_DB_CONNECTIONS", "100")
	setEnv(t, "PRINT_SCRAPE_ScreenShot_STORAGE_FOLDER", "./")

	defer unsetEnv(t, "PRINT_SCRAPE_DB_HOST")
	defer unsetEnv(t, "PRINT_SCRAPE_DB_NAME")
	defer unsetEnv(t, "PRINT_SCRAPE_DB_PORT")
	defer unsetEnv(t, "PRINT_SCRAPE_DB_USER")
	defer unsetEnv(t, "PRINT_SCRAPE_DB_PASSWORD")
	defer unsetEnv(t, "MAX_DB_CONNECTIONS")
	defer unsetEnv(t, "PRINT_SCRAPE_ScreenShot_STORAGE_FOLDER")

	cfg := config.FromEnv()

	if got, want := cfg.DatabaseUser, "dbuser"; got != want {
		t.Errorf("env var PRINT_SCRAPE_DB_USER=%q, want %q", got, want)
	}
	if got, want := cfg.DatabaseHost, "host"; got != want {
		t.Errorf("env var PRINT_SCRAPE_DB_HOST=%q,want %q", got, want)
	}
	if got, want := cfg.DatabaseName, "db"; got != want {
		t.Errorf("env var PRINT_SCRAPE_DB_NAME=%q,want %q", got, want)
	}
	if got, want := cfg.DatabasePort, "5000"; got != want {
		t.Errorf("env var PRINT_SCRAPE_DB_PORT=%q,want %q", got, want)
	}
	if got, want := cfg.DatabasePassword, "password"; got != want {
		t.Errorf("env var PRINT_SCRAPE_DB_PASSWORD=%q,want %q", got, want)
	}

	if got, want := cfg.MaxDBConnections, 100; got != want {
		t.Errorf("env var MAX_DB_CONNECTIONS=%d,want %d", got, want)
	}

	if got, want := cfg.ScreenShotStorageFolder, "./"; got != want {
		t.Errorf("env var PRINT_SCRAPE_ScreenShot_STORAGE_FOLDER=%s,want %s", got, want)
	}
}

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed setting env %q as %q: %v", key, value, err)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed unsetting env %q: %v", key, err)
	}
}
