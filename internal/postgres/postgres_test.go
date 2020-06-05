package postgres_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
	"github.com/slysterous/print-scrape/internal/postgres"
	"testing"
)

func TestGetScrapSuccess(t *testing.T) {
	db, mock, closeDB := sqlMockNew(t)

	defer closeDB(db)

	client := postgres.Client{DB: db}

	want := printscrape.Screenshot{
		FileURI: "testfile",
		RefCode: "testcode",
	}

	const query = "SELECT fileURI FROM screenshots WHERE refCode\\=.*"

	columns := []string{"FileURI"}

	mock.ExpectQuery(query).WillReturnRows(mock.NewRows(columns).AddRow(
		want.FileURI,
	))

	got, err := client.GetScrapByCode("testcode")
	if err != nil {
		t.Fatalf("expected exec not to return error %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	if got.FileURI != want.FileURI {
		t.Fatalf("expected: %s, got: %s", want.FileURI, got.FileURI)
	}

	if want.RefCode != want.RefCode {
		t.Fatalf("expected: %s, got: %s", want.RefCode, got.RefCode)
	}
}

func TestCreateScrapSuccess(t *testing.T) {
	db, mock, closeDB := sqlMockNew(t)

	defer closeDB(db)

	client := postgres.Client{DB: db}

	want := printscrape.Screenshot{
		FileURI: "testfile",
		RefCode: "testcode",
	}

	const query = "INSERT INTO scraps \\(.*\\) VALUES .* RETURNING id"
	mock.ExpectQuery(query).WithArgs(want.RefCode, want.FileURI).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	got, err := client.CreateScrap(want)
	if err != nil {
		t.Fatalf("could not create scrap: %v", err)
	}
	if got != 1 {
		t.Errorf("got=%d, want=%d", got, 1)
	}
}

// sqlMockNew mock sql connection.
func sqlMockNew(t *testing.T) (db *sql.DB, mock sqlmock.Sqlmock, closeDB func(db *sql.DB)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a stub database connection", err)
	}

	closeDB = func(db *sql.DB) {
		_ = db.Close()
	}
	return db, mock, closeDB
}
