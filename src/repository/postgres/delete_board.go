package postgres

import (
	"context"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/usecase/deleteboard"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Repository) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	const op = "postgres.DeleteBoard"
	exists := r.CheckBoard(ctx, id.String())
	if !exists {
		return deleteboard.ErrBoardDoesntExist
	}

	now := time.Now().UTC()
	_, err := r.pool.Exec(ctx, `
	UPDATE boards SET 
    deleted_at = $1,
    updated_at = $2 
	WHERE id = $3 AND deleted_at IS NULL`, now, now, id)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
