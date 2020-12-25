package postgres_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/slysterous/scrapmon/internal/postgres"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
	"testing"
	"time"
)

func TestGetLatestCreatedScrapCode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := "150000"

		const query = "SELECT refCode from Scraps ORDER BY codeCreatedAt DESC limit 1"

		columns := []string{"RefCode"}

		mock.ExpectQuery(query).WillReturnRows(mock.NewRows(columns).AddRow(
			want,
		))

		got, err := client.GetLatestCreatedScrapCode()
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
	})
	t.Run("Error", func (t *testing.T){
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "SELECT refCode from Scraps ORDER BY codeCreatedAt DESC limit 1"

		mock.ExpectQuery(query).WillReturnError(errors.New("test error"))

		got, err := client.GetLatestCreatedScrapCode()
		if err == nil {
			t.Fatalf("expected err, got:%v", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}

		if got != nil {
			t.Fatalf("expected: nil, got: %s", *got)
		}
	})
	t.Run("No Codes", func (t *testing.T){
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "SELECT refCode from Scraps ORDER BY codeCreatedAt DESC limit 1"

		columns := []string{"RefCode"}

		mock.ExpectQuery(query).WillReturnRows(mock.NewRows(columns))

		got, err := client.GetLatestCreatedScrapCode()
		if err != nil {
			t.Fatalf("expected exec not to return error %v", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}

		if got !=nil {
			t.Fatalf("expected: nil, got: %s", *got)
		}
	})
}

func TestGetScrap(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := scrapmon.Scrap{
			FileURI: "testfile",
			RefCode: "testcode",
		}

		const query = "SELECT fileURI FROM Scraps WHERE refCode\\=.*"

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
	})
	t.Run("Error", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}
		const query = "SELECT fileURI FROM Scraps WHERE refCode\\=.*"

		mock.ExpectQuery(query).WillReturnError(errors.New(""))

		_, err := client.GetScrapByCode("testcode")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}
	})
	t.Run("Not Found", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}
		const query = "SELECT fileURI FROM Scraps WHERE refCode\\=.*"

		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		got, err := client.GetScrapByCode("testcode")
		if err != nil {
			t.Fatalf("expected nil, got error: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil, got: %v", got)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}
	})
}

func TestCreateScrap(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := scrapmon.Scrap{
			FileURI:       "testfile",
			RefCode:       "testcode",
			CodeCreatedAt: time.Now(),
		}

		const query = "INSERT INTO Scraps \\(.*\\) VALUES .* RETURNING id"
		mock.ExpectQuery(query).WithArgs(want.RefCode, want.CodeCreatedAt, want.FileURI).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		got, err := client.CreateScrap(want)
		if err != nil {
			t.Fatalf("could not create scrap: %v", err)
		}
		if got != 1 {
			t.Errorf("got=%d, want=%d", got, 1)
		}
	})
	t.Run("Error", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := scrapmon.Scrap{
			FileURI:       "testfile",
			RefCode:       "testcode",
			CodeCreatedAt: time.Now(),
		}

		const query = "INSERT INTO Scraps \\(.*\\) VALUES .* RETURNING id"
		mock.ExpectQuery(query).WithArgs(want.RefCode, want.CodeCreatedAt, want.FileURI).WillReturnError(errors.New("test error"))

		_, err := client.CreateScrap(want)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUpdateScrapStatusByCode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}
		refCode := "testcode"

		const query = "UPDATE Scraps SET downloadStatus = .* WHERE refCode = .*"
		mock.ExpectExec(query).WithArgs(scrapmon.StatusSuccess, refCode).WillReturnResult(sqlmock.NewResult(0, 1))

		err := client.UpdateScrapStatusByCode(refCode, scrapmon.StatusSuccess)
		if err != nil {
			t.Fatalf("expected nil, got err: %v", err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}
		refCode := "testcode"

		const query = "UPDATE Scraps SET downloadStatus = .* WHERE refCode = .*"
		mock.ExpectExec(query).WithArgs(scrapmon.StatusSuccess, refCode).WillReturnError(errors.New("test error"))

		err := client.UpdateScrapStatusByCode(refCode, scrapmon.StatusSuccess)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

func TestUpdateScrapByCode(t *testing.T) {
	t.Run("Success", func(t *testing.T){
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := scrapmon.Scrap{
			FileURI:       "testfile",
			RefCode:       "testcode",
			Status: scrapmon.StatusFailure,
		}

		const query = "UPDATE Scraps SET fileUri= .*,downloadStatus= .* WHERE refCode = .*;"
		mock.ExpectExec(query).WithArgs(want.FileURI, scrapmon.StatusFailure, want.RefCode).WillReturnResult(sqlmock.NewResult(0, 1))

		err := client.UpdateScrapByCode(want)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})
	t.Run("Error", func(t *testing.T){
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		want := scrapmon.Scrap{
			FileURI:       "testfile",
			RefCode:       "testcode",
			Status: scrapmon.StatusFailure,
		}

		const query = "UPDATE Scraps SET fileUri= .*,downloadStatus= .* WHERE refCode = .*;"
		mock.ExpectExec(query).WithArgs(want.FileURI, scrapmon.StatusFailure, want.RefCode).WillReturnError(errors.New("test error"))
		err := client.UpdateScrapByCode(want)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCodeAlreadyExists(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "SELECT id FROM Scraps WHERE refCode\\=.*"

		columns := []string{"id"}

		mock.ExpectQuery(query).WillReturnRows(mock.NewRows(columns).AddRow(
			"1",
		))

		got, err := client.CodeAlreadyExists("testcode")
		if err != nil {
			t.Fatalf("expected query not to return error, got: %v", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}

		if got != true {
			t.Fatalf("expected: %v, got: %v", true, got)
		}
	})
	t.Run("Error", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "SELECT id FROM Scraps WHERE refCode\\=.*"

		mock.ExpectQuery(query).WillReturnError(errors.New("test error"))

		got, err := client.CodeAlreadyExists("testcode")
		if err == nil {
			t.Fatal("expected error got nil")
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}

		if got != false {
			t.Fatalf("expected false, got: %v", got)
		}
	})
	t.Run("Not Found", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "SELECT id FROM Scraps WHERE refCode\\=.*"

		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		got, err := client.CodeAlreadyExists("testcode")

		if err != nil {
			t.Fatalf("expected nil, got error: %v", err)
		}
		if got != false {
			t.Fatalf("expected false, got: %v", got)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}
	})
}

func TestPurge(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "TRUNCATE TABLE Scraps"

		mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))

		err := client.Purge()
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}
	})
	t.Run("Error", func(t *testing.T) {
		db, mock, closeDB := sqlMockNew(t)

		defer closeDB(db)

		client := postgres.Client{DB: db}

		const query = "TRUNCATE TABLE Scraps"

		mock.ExpectExec(query).WillReturnError(errors.New("test error"))

		err := client.Purge()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}
	})
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
