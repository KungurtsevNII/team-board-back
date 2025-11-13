package postgres

import (
	"context"
	"errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
)

func (r Repository) GetBoards(user_id uuid.UUID, ctx context.Context) (*[]domain.Board, error) {
	//TODO проверять user_id и подставлять его как параметр в запрос
	_, err := uuid.Parse(user_id.String())
	if err != nil {
		return nil, ErrInvalidUserID
	}

	boards := make([]domain.Board, 0)
	//TODO подставлять user_id в запрос
	rows, err := r.pool.Query(ctx, `
	SELECT id, name, short_name, created_at, updated_at, deleted_at
	FROM boards
	ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var board domain.Board
		err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.ShortName,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return &boards, nil
}

var ErrInvalidUserID = errors.New("invalid user id")
