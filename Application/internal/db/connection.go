package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var SqlOpen = sql.Open
var SqlPing = func(db *sql.DB) error {
	return db.Ping()
}

type PostgresDB struct {
	db DatabaseInterface
}

func (p *PostgresDB) Query(query string, args ...any) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *PostgresDB) QueryRow(query string, args ...any) *sql.Row {
	return p.db.QueryRow(query, args...)
}

func (p *PostgresDB) Exec(query string, args ...any) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// ConnectDB returns a DatabaseInterface without using a global.
func ConnectDB(connStr string) (DatabaseInterface, error) {
	db, err := SqlOpen("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %v", err)
	}

	if err := SqlPing(db); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return &PostgresDB{db: db}, nil
}

func NewPostgresDB(mock DatabaseInterface) *PostgresDB {
	return &PostgresDB{db: mock}
}
