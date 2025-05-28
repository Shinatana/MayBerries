package init

import (
	"auth_service/internal/repo/migrate"
	"auth_service/internal/repo/pgsql"
	"auth_service/pkg/config"
	"context"
	"fmt"
)

func init_db(dbOptions *config.DatabaseOptions, migrationOptions *config.MigrationOptions) (*pgsql.Pgsql, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOptions.InitTimeout)
	defer cancel()

	err := migrate.NewMigrator(dbOptions, migrationOptions).Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	db, err := pgsql.NewDB(ctx, dbOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	return db, nil
}
