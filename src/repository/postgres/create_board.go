package postgres

import (
	"context"

	"github.com/pkg/errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
)

func (r Repository) CreateBoard(ctx context.Context, board domain.Board) error {
	op := "postgres.CreateBoard"

	_, err := r.pool.Exec(ctx,
		`INSERT INTO boards (id, name, short_name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		board.ID,
		board.Name,
		board.ShortName,
		board.CreatedAt,
		board.UpdatedAt,
		board.DeletedAt,
	)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
