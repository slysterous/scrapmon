package test

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	migPostgres "github.com/golang-migrate/migrate/database/postgres"
	"github.com/slysterous/print-scrape/internal/postgres"

	//Bilateral import
	_ "github.com/golang-migrate/migrate/source/file"
)

// maxConnections defines the max database connections.
const maxDBConnections = 5

// TearUp sets up the database to be used in testing.
func TearUp(t *testing.T) *sql.DB {
	db, err := postgres.NewClient(getDataSource(), maxDBConnections)
	if err != nil {
		t.Fatalf("could not connect to dabase: %v", err)
	}

	err = runMigrations(t, db.DB)
	if err != nil {
		t.Fatalf("could not run migrations: %v", err)
	}

	err = truncateTables(db.DB)
	if err != nil {
		t.Fatalf("could not truncate tables: %v", err)
	}

	return db.DB
}

// TearDown closes the database connection.
func TearDown(db io.Closer, t *testing.T) {
	err := db.Close()
	if err != nil {
		t.Fatalf("could not close database connection %v", err)
	}
}

func getDataSource() string {
	user := getEnv("PRINT_SCRAPE_DB_USER", "postgres")
	pass := getEnv("PRINT_SCRAPE_DB_PASSWORD", "password")
	host := getEnv("PRINT_SCRAPE_DB_HOST", "127.0.0.1")
	port := getEnv("PRINT_SCRAPE_DB_PORT", "5432")
	name := getEnv("DB_NAME", "print-scrape")
	return "host=" + host + " port=" + port + " user=" + user + " password=" + pass + " dbname=" + name + " sslmode=disable"
}

// getEnv will return the environment variable or the default value.
func getEnv(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

// truncateTables truncates all the tables of the database.
// The migrations table is skipped.
func truncateTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Disable FK checks, add truncate table commands and then reset FK checks
	cmds := []string{"SET session_replication_role = 'replica'"}
	tables := []string{
		"scraps",
	}
	for _, table := range tables {
		cmds = append(cmds, fmt.Sprintf("TRUNCATE %s CASCADE", table))
	}
	cmds = append(cmds, "SET session_replication_role = 'origin'")

	// Perform all checks in a single transaction and revert if anything goes wrong
	for _, cmd := range cmds {
		if _, cmdErr := tx.Exec(cmd); cmdErr != nil {
			_ = tx.Rollback()
			return cmdErr
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func runMigrations(t *testing.T, db *sql.DB) error {
	if ok, err := shouldRunMigrations(db); !ok || err != nil {
		return err
	}

	migDriver, err := migPostgres.WithInstance(db, &migPostgres.Config{})
	if err != nil {
		t.Fatalf("could not create migration driver %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations",
		"postgres",
		migDriver,
	)
	if err != nil {
		t.Fatalf("could not create migration database instance: %v", err)
	}

	err = m.Up()
	if err != nil {
		t.Fatalf("could not run migrations: %v", err)
	}
	return nil
}

func shouldRunMigrations(db *sql.DB) (bool, error) {
	var numOfTables int
	row := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
	err := row.Scan(&numOfTables)
	if err != nil {
		return false, err
	}
	return numOfTables == 0, nil
}
