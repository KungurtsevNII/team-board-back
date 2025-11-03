package postgres

import (
	"database/sql"
	"fmt"
)


type Repository struct {
	db *sql.DB
}

func New(storagePath string) (*Repository, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{db: db}, nil
}

func (s *Repository) Close() error {
	op := "storage.postgresql.Close"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}