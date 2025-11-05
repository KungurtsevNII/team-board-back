package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(storagePath string) (*Postgres, error) {
	const op = "storage.postgresql.New"

	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{
		pool: pool,
	}, nil
}

func (s *Postgres) Close() {
	s.pool.Close()
}
