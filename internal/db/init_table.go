package db

import (
	"database/sql"
	"fmt"
	"os"
)

// ExecuteSQLFile reads an SQL file and executes its statements on the given db
func ExecuteSQLFile(db *sql.DB, filePath string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	sqlStmt := string(sqlBytes)

	// Execute all SQL statements
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to execute SQL file: %w", err)
	}

	return nil
}
