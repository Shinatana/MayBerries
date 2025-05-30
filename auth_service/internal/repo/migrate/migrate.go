package migrations

import (
	"database/sql"
	"fmt"

	"auth_service/pkg/log"
)

func Up(tx *sql.Tx) error {
	log.Info("⏫ Starting migration: create auth schema")

	queries := []string{
		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			description TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			code TEXT UNIQUE NOT NULL,
			description TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT NOT NULL,
			role_id INTEGER REFERENCES roles(id) ON DELETE SET NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS roles_permissions (
			role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
			permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, permission_id)
		);`,
	}

	for _, query := range queries {
		log.Debug("Executing migration query", "sql", query)
		if _, err := tx.Exec(query); err != nil {
			log.Error("Failed migration query", "sql", query, "error", err)
			return fmt.Errorf("migration failed on query: %s\n%w", query, err)
		}
	}

	log.Info("✅ Migration complete: auth schema created")
	return nil
}

func Down(tx *sql.Tx) error {
	log.Info("⏬ Starting rollback: drop auth schema")

	queries := []string{
		`DROP TABLE IF EXISTS roles_permissions;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS permissions;`,
		`DROP TABLE IF EXISTS roles;`,
	}

	for _, query := range queries {
		log.Debug("Executing rollback query", "sql", query)
		if _, err := tx.Exec(query); err != nil {
			log.Error("Failed rollback query", "sql", query, "error", err)
			return fmt.Errorf("rollback failed on query: %s\n%w", query, err)
		}
	}

	log.Info("✅ Rollback complete: auth schema dropped")
	return nil
}
