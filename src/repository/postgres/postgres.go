package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Repository struct {
	pool *pgxpool.Pool
}

type RepositoryInf interface {
    
}

func New(storagePath string) (*Repository, error) {
	const op = "storage.postgresql.New"

	log := slog.Default()
	log.With("storagePath", storagePath, "op", op).Info("connecting to postgres")
	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (s *Repository) Close(){
	s.pool.Close()
}
