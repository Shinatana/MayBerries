package repo

import "context"

// Migrator defines an abstraction for running migrations.
type Migrator interface {
	Migrate(ctx context.Context) error
}

type DB interface {
	Close()
}
