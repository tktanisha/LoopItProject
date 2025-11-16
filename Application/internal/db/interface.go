package db

import (
	"database/sql"
)

//go:generate mockgen -source=interface.go -destination=../mock/mock_db.go -package=mock

type DatabaseInterface interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	Close() error
}

type RealDatabase struct {
	*sql.DB
}

func NewRealDatabase(db *sql.DB) DatabaseInterface {
	return &RealDatabase{DB: db}
}
