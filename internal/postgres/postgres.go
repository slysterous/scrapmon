package postgres

import (
	"database/sql"
	"fmt"
	"time"
	//pq is the postgres driver for database/sql
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	printscrape "github.com/slysterous/print-scrape/internal/printscrape"
)

// Client represents the postgres client.
type Client struct {
	DB *sql.DB
}

type scrap struct {
	id             int64
	uuid           uuid.UUID
	downloadStatus string
	refCode        string
	codeCreatedAt  time.Time
	fileURI        string
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

// CreateScreenShot saves a Scrap's info in the database.
func (c *Client) CreateScreenShot(ss printscrape.ScreenShot) (int, error) {
	scrap := transformDomainScrap(ss)
	const query = `INSERT INTO screenshots (refCode,codeCreatedAt,fileUri) VALUES ($1, $2, $3) RETURNING id`
	lastInsertID := 0

	err := c.DB.QueryRow(query, scrap.refCode, scrap.codeCreatedAt, scrap.fileURI).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, fmt.Errorf("postgres: executing insert screenshot statement: %v", err)
	}
	return lastInsertID, nil
}

func (c *Client) UpdateScreenShotStatusByCode(code string, status printscrape.ScreenShotStatus) error {
	const query = `UPDATE screenshots SET downloadStatus = $1 WHERE refCode = $2;`
	_, err := c.DB.Exec(query, string(status), code)
	if err != nil {
		return fmt.Errorf("postgres: executing update status on scrap: %v", err)
	}
	return nil
}

func (c *Client) UpdateScreenShotByCode(ss printscrape.ScreenShot) error {
	scrap := transformDomainScrap(ss)
	const query = `UPDATE screenshots SET fileUri = $1,downloadStatus = $2 WHERE refCode = $3;`
	_, err := c.DB.Exec(query, scrap.fileURI, scrap.downloadStatus, scrap.refCode)
	if err != nil {
		return fmt.Errorf("postgres: executing update of scrap: %v", err)
	}
	return nil
}

func (c *Client) GetLatestCreatedScreenShotCode() (*string, error) {
	const query = `SELECT refCode from screenshots ORDER BY codeCreatedAt DESC limit 1`
	code := ""
	row := c.DB.QueryRow(query)

	err := row.Scan(&code)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		fmt.Println("GAMIOLA")
		return &code, nil
	default:
		return nil, fmt.Errorf("postgres: executing select screenshot by code statement: %v", err)
	}

}

// CodeAlreadyExists searches for the existence of an entry.
func (c *Client) CodeAlreadyExists(code string) (bool, error) {
	s := scrap{
		id:      0,
		uuid:    uuid.UUID{},
		refCode: code,
		fileURI: "",
	}
	row := c.DB.QueryRow(`
		SELECT
		id
		FROM ScreenShots
		WHERE refCode=$1
	
	`, code)

	err := row.Scan(&s.fileURI)

	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("postgres: executing image already exists statement: %v", err)
	}
}

// GetScrapByCode Gets a scrapped ScreenShot from the database
func (c *Client) GetScrapByCode(code string) (*printscrape.ScreenShot, error) {
	s := scrap{
		id:      0,
		uuid:    uuid.UUID{},
		refCode: code,
		fileURI: "",
	}
	row := c.DB.QueryRow(`
		SELECT
		fileURI
		FROM ScreenShots
		WHERE refCode=$1
	
	`, code)

	err := row.Scan(&s.fileURI)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return transformScrap(s), nil
	default:
		return nil, fmt.Errorf("postgres: executing select scrap by code statement: %v", err)
	}
}

func (c *Client) Purge() error {
	if _, err := c.DB.Exec(`TRUNCATE TABLE screenshots`); err != nil {
		return fmt.Errorf("postgres: executing truncate query, err: %v", err)
	}
	return nil
}

// transformScrap transforms the postgres scrap to a domain Scrap
func transformScrap(sc scrap) *printscrape.ScreenShot {
	return &printscrape.ScreenShot{
		RefCode: sc.refCode,
		FileURI: sc.fileURI,
	}
}

func transformDomainScrap(screenShot printscrape.ScreenShot) scrap {
	return scrap{
		id:             screenShot.ID,
		codeCreatedAt:  screenShot.CodeCreatedAt,
		downloadStatus: string(screenShot.Status),
		refCode:        screenShot.RefCode,
		fileURI:        screenShot.FileURI,
	}
}
