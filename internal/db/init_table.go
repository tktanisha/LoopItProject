package db

import (
	"database/sql"
	"fmt"
	"os"
)

func ExecuteSQLFile(db *sql.DB, filePath string) error {

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	sqlStmt := string(sqlBytes)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file: %w", err)
	}

	return nil
}
