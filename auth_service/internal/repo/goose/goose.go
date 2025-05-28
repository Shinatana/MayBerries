package migrations

import (
	"database/sql"
	"fmt"
)

// Up выполняется при применении миграции
func Up(tx *sql.Tx) error {
	query := `UPDATE users SET username = 'admin' WHERE username = 'root';`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to update username from 'root' to 'admin': %w", err)
	}
	return nil
}

// Down выполняется при откате миграции
func Down(tx *sql.Tx) error {
	query := `UPDATE users SET username = 'root' WHERE username = 'admin';`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("failed to revert username from 'admin' to 'root': %w", err)
	}
	return nil
}
