package postgres

import (
	"context"
	"errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r Repository) GetBoard(ctx context.Context, ID string) (domain.Board, error) {
	uid, err := uuid.Parse(ID)
	if err != nil {
		return domain.Board{}, domain.ErrInvalidID
	}

	var board domain.Board
	err = pgxscan.Get(ctx, r.pool, &board,
		`SELECT id, name, short_name, created_at, updated_at, deleted_at 
		FROM boards WHERE id = $1`, uid)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Board{}, domain.ErrBoardNotFound
		}
		return domain.Board{}, err
	}

	board.Columns = make([]domain.Column, 0)
	err = pgxscan.Select(ctx, r.pool, &board.Columns,
		`SELECT id, board_id, order_num, name 
		FROM columns WHERE board_id = $1 
		ORDER BY order_num;`, uid)
	if err != nil {
		return domain.Board{}, err
	}

	board.Tasks = make([]domain.Task, 0)
	err = pgxscan.Select(ctx, r.pool, &board.Tasks,
		`SELECT id, column_id, board_id, number, title 
		FROM tasks WHERE board_id = $1 
		ORDER BY number;`, uid)
	if err != nil {
		return domain.Board{}, err
	}

	return board, nil

}
