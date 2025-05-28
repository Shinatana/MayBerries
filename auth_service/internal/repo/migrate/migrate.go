package migrate

import (
	"auth_service/pkg/config"
	"auth_service/pkg/log"
	"auth_service/pkg/misc"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
)

type migratorImpl struct {
	MigrationFiles string
	Version        int
	DSN            string
}

func NewMigrator(db *config.DatabaseOptions, migration *config.MigrationOptions) *migratorImpl {
	return &migratorImpl{
		MigrationFiles: migration.MigrationFiles,
		DSN:            misc.GetDSN(db, misc.WithMigratorFormat()), // без замены на pgx5
		Version:        migration.Version,
	}
}

func (m *migratorImpl) Migrate(ctx context.Context) error {
	if m.Version == 0 {
		log.Info("migrations skipped due to configuration settings", "version", m.Version)
		return nil
	}

	errC := make(chan error, 1)

	go func() {
		defer close(errC)
		errC <- m.runMigration(ctx)
	}()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
		select {
		case err := <-errC:
			return err
		default:
			return ctx.Err()
		}
	}
}

func (m *migratorImpl) runMigration(ctx context.Context) error {
	migrateInstance, err := migrate.New("file://"+m.MigrationFiles, m.DSN)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() { _, _ = migrateInstance.Close() }()

	if err = ctx.Err(); err != nil {
		return err
	}

	log.Debug("migration files path", "path", m.MigrationFiles)
	version, dirty, err := migrateInstance.Version()
	log.Info("migration status", "version", version, "dirty", dirty, "error", err)

	if err = ctx.Err(); err != nil {
		return err
	}

	if m.Version > 0 {
		if err = migrateInstance.Migrate(uint(m.Version)); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate to version %d: %w", m.Version, err)
		} else if errors.Is(err, migrate.ErrNoChange) {
			log.Info("already at requested version", "version", m.Version)
		} else {
			log.Info("migrated to specific version", "version", m.Version)
		}
	} else {
		if err = migrateInstance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		} else if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migration required")
		}
	}

	if version1, dirty1, err1 := migrateInstance.Version(); version != version1 {
		log.Info("new migration status", "version", version1, "dirty", dirty1, "error", err1)
		log.Info("migration successful")
	}

	return nil
}
