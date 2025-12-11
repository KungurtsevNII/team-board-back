package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetBoard(ctx context.Context, ID uuid.UUID) (*domain.Board, error) {
	const op = "postgres.GetBoard"

	var board domain.Board
	err := pgxscan.Get(ctx, r.pool, &board,
		`SELECT id, name, short_name, created_at, updated_at, deleted_at 
		FROM boards WHERE id = $1
		AND deleted_at IS NULL`, ID)

	if err != nil {
		return nil, err
	}

	board.Columns, err = r.GetColumns(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	board.Tasks, err = r.GetTasks(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &board, nil

}
