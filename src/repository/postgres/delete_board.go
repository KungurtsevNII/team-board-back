package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Repository) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	const op = "postgres.DeleteBoard"

	_, err := r.pool.Exec(ctx, `
	UPDATE boards SET 
    deleted_at = NOW(),
    updated_at = NOW() 
	WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
