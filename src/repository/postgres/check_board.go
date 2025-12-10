package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (r Repository) CheckBoard(ctx context.Context, id string) bool {
	uid, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	var exists bool
	err = r.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM boards WHERE id = $1 AND deleted_at IS NULL)`,
		uid,
	).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
