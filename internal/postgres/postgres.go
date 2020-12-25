package postgres

import (
	"database/sql"
	"fmt"
	"time"
	//pq is the postgres driver for database/sql
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	scrapmon "github.com/slysterous/scrapmon/internal/scrapmon"
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
	//TODO where is the close???
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
func (c *Client) CreateScrap(ss scrapmon.Scrap) (int, error) {
	scrap := transformDomainScrap(ss)
	const query = `INSERT INTO Scraps (refCode,codeCreatedAt,fileUri) VALUES ($1, $2, $3) RETURNING id`
	lastInsertID := 0

	err := c.DB.QueryRow(query, scrap.refCode, scrap.codeCreatedAt, scrap.fileURI).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, fmt.Errorf("postgres: executing insert Scrap statement: %v", err)
	}
	return lastInsertID, nil
}

func (c *Client) UpdateScrapStatusByCode(code string, status scrapmon.ScrapStatus) error {
	const query = `UPDATE Scraps SET downloadStatus = $1 WHERE refCode = $2;`
	_, err := c.DB.Exec(query, string(status), code)
	if err != nil {
		return fmt.Errorf("postgres: executing update status on scrap: %v", err)
	}
	return nil
}

func (c *Client) UpdateScrapByCode(ss scrapmon.Scrap) error {
	scrap := transformDomainScrap(ss)
	const query = `UPDATE Scraps SET fileUri= $1,downloadStatus= $2 WHERE refCode = $3;`
	_, err := c.DB.Exec(query, scrap.fileURI, scrap.downloadStatus, scrap.refCode)
	if err != nil {
		return fmt.Errorf("postgres: executing update of scrap: %v", err)
	}
	return nil
}

func (c *Client) GetLatestCreatedScrapCode() (*string, error) {
	const query = `SELECT refCode from Scraps ORDER BY codeCreatedAt DESC limit 1`
	code := ""
	row := c.DB.QueryRow(query)

	err := row.Scan(&code)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &code, nil
	default:
		return nil, fmt.Errorf("postgres: executing select Scrap by code statement: %v", err)
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
		FROM Scraps
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

// GetScrapByCode Gets a scrapped Scrap from the database
func (c *Client) GetScrapByCode(code string) (*scrapmon.Scrap, error) {
	s := scrap{
		id:      0,
		uuid:    uuid.UUID{},
		refCode: code,
		fileURI: "",
	}
	row := c.DB.QueryRow(`
		SELECT
		fileURI
		FROM Scraps
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
	if _, err := c.DB.Exec(`TRUNCATE TABLE Scraps`); err != nil {
		return fmt.Errorf("postgres: executing truncate query, err: %v", err)
	}
	return nil
}

// transformScrap transforms the postgres scrap to a domain Scrap
func transformScrap(sc scrap) *scrapmon.Scrap {
	return &scrapmon.Scrap{
		RefCode: sc.refCode,
		FileURI: sc.fileURI,
	}
}

func transformDomainScrap(Scrap scrapmon.Scrap) scrap {
	return scrap{
		id:             Scrap.ID,
		codeCreatedAt:  Scrap.CodeCreatedAt,
		downloadStatus: string(Scrap.Status),
		refCode:        Scrap.RefCode,
		fileURI:        Scrap.FileURI,
	}
}
