package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
)

func (r Repository) CreateBoard(board domain.Board, ctx context.Context) error {
	const sql = `
		INSERT INTO boards (id, name, short_name, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, sql,
		board.ID,
		board.Name,
		board.ShortName,
		board.CreatedAt,
		board.UpdatedAt,
		board.DeletedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
