package postgres_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/slysterous/print-scrape/internal/postgres"
	printscrape "github.com/slysterous/print-scrape/internal/printscrape"
	"testing"
	"time"
)

func TestGetLatestCreatedScrapCodeSuccess(t *testing.T) {
	db, mock, closeDB := sqlMockNew(t)

	defer closeDB(db)

	client := postgres.Client{DB: db}

	want := "150000"

	const query = "SELECT refCode from screenshots ORDER BY codeCreatedAt DESC limit 1"

	columns := []string{"RefCode"}

	mock.ExpectQuery(query).WillReturnRows(mock.NewRows(columns).AddRow(
		want,
	))

	got, err := client.GetLatestCreatedScreenShotCode()
	if err != nil {
		t.Fatalf("expected exec not to return error %v", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	if want != *got {
		t.Fatalf("expected: %s, got: %s", want, *got)
	}

}

func TestGetScrapSuccess(t *testing.T) {
	db, mock, closeDB := sqlMockNew(t)

	defer closeDB(db)

	client := postgres.Client{DB: db}

	want := printscrape.ScreenShot{
		FileURI: "testfile",
		RefCode: "testcode",
	}

	const query = "SELECT fileURI FROM ScreenShots WHERE refCode\\=.*"

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

	want := printscrape.ScreenShot{
		FileURI:       "testfile",
		RefCode:       "testcode",
		CodeCreatedAt: time.Now(),
	}

	const query = "INSERT INTO screenshots \\(.*\\) VALUES .* RETURNING id"
	mock.ExpectQuery(query).WithArgs(want.RefCode, want.CodeCreatedAt, want.FileURI).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	got, err := client.CreateScreenShot(want)
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
