package pgsql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"auth_service/internal/repo"
	"auth_service/pkg/config"
	"auth_service/pkg/misc"
)

const initPingTimeout = 1 * time.Second

type pgsql struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dbConfig *config.DatabaseOptions) (repo.DB, error) {
	config, err := pgxpool.ParseConfig(misc.GetDSN(dbConfig, misc.WithPGXv5Format()))
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	ctsSec, cancel := context.WithTimeout(ctx, initPingTimeout)
	defer cancel()

	err = pool.Ping(ctsSec)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping pool: %w", err)
	}

	return &pgsql{pool: pool}, nil
}

func (p *pgsql) Close() {
	p.pool.Close()
}
