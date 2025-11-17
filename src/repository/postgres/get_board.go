package postgres

import (
	"context"
	"errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r Repository) GetBoard(ctx context.Context, ID string) (domain.Board, error) {
	uid, err := uuid.Parse(ID)
	if err != nil {
		return domain.Board{}, domain.ErrInvalidID
	}
	var board domain.Board

	//тут получаем доску
	row := r.pool.QueryRow(ctx,
		`SELECT id, name, short_name, created_at, updated_at, deleted_at 
		FROM boards WHERE id = $1`,
		uid,
	)

	err = row.Scan(
		&board.ID,
		&board.Name,
		&board.ShortName,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.Board{}, domain.ErrBoardNotFound
		}
		return domain.Board{}, err
	}

	//тут получаем колонки доски
	board.Columns = make([]domain.Column, 0)
	rows, err := r.pool.Query(ctx, `
	SELECT id, board_id, order_num, name 
	FROM columns WHERE board_id = $1 
	ORDER BY order_num;`, uid)
	if err != nil {
		return domain.Board{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var column domain.Column
		err := rows.Scan(
			&column.ID,
			&column.BoardID,
			&column.OrderNum,
			&column.Name,
		)
		if err != nil {
			return domain.Board{}, err
		}
		board.Columns = append(board.Columns, column)
	}

	//TODO прикрутить задачи в колонки , сделаем как замерджим структуру task , чтобы не путать ветки

	return board, nil
}
