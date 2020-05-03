package postgres

import (
	"database/sql"
	"fmt"
)

// Client represents the postgres client.
type Client struct {
	DB *sql.DB
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
