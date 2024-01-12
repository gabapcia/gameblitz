package postgres

import (
	"context"

	"github.com/gabarcia/metagaming-api/internal/infra/storage/postgres/internal/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connection struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func (c connection) Close() {
	c.pool.Close()
}

func New(ctx context.Context, connStr string) (*connection, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	conn := &connection{
		pool:    pool,
		queries: sqlc.New(pool),
	}

	return conn, nil
}
