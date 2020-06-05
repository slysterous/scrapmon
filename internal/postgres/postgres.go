package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	printscrape "github.com/slysterous/print-scrape/internal/domain"
)

// Client represents the postgres client.
type Client struct {
	DB *sql.DB
}

type scrap struct {
	id      int64
	uuid    uuid.UUID
	refCode string
	fileURI string
}

// NewClient returns a postgres database client to interact with.
func NewClient(dataSource string, maxConnections int) (*Client, error) {
	conn, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("postgres: connecting to database: %v", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres: pinging the database: %v", err)
	}

	conn.SetMaxOpenConns(maxConnections)

	client := &Client{
		DB: conn,
	}
	return client, nil
}

// CreateScrap saves a Scrap's info in the database.
func (c *Client) CreateScrap(scrap printscrape.Screenshot) (int, error) {
	query := `INSERT INTO scraps (refCode, fileUri) VALUES ($1, $2) RETURNING id`
	lastInsertID := 0

	err := c.DB.QueryRow(query, scrap.RefCode, scrap.FileURI).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, fmt.Errorf("postgres: executing insert scrap statement: %v", err)
	}
	return lastInsertID, nil
}

// GetScrapByCode Gets a scrapped screenshot from the database
func (c *Client) GetScrapByCode(code string) (sc *printscrape.Screenshot, err error) {
	s := scrap{
		id:      0,
		uuid:    uuid.UUID{},
		refCode: code,
		fileURI: "",
	}
	row := c.DB.QueryRow(`
		SELECT
		fileURI
		FROM screenshots
		WHERE refCode=$1
	
	`, code)

	err = row.Scan(&s.fileURI)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return transformScrap(s), nil
	default:
		return nil, fmt.Errorf("postgres: executing select scrap by code statement: %v", err)
	}
}

// transformScrap transforms the postgres scrap to a domain Scrap
func transformScrap(sc scrap) *printscrape.Screenshot {
	return &printscrape.Screenshot{
		RefCode: sc.refCode,
		FileURI: sc.fileURI,
	}
}
