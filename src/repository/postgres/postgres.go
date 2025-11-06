package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

type RepositoryInf interface {
    
}

func New(storagePath string) (*Repository, error) {
	const op = "storage.postgresql.New"

	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (s *Repository) Close(){
	s.pool.Close()
}
