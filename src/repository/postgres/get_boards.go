package postgres

import (
	"context"
	"errors"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

func (r Repository) GetBoards(user_id uuid.UUID, ctx context.Context) ([]domain.Board, error) {
	//TODO проверять user_id и подставлять его как параметр в запрос
	_, err := uuid.Parse(user_id.String())
	if err != nil {
		return nil, ErrInvalidUserID
	}

	boards := make([]domain.Board, 0)
	//TODO подставлять user_id в запрос
	err = pgxscan.Select(ctx, r.pool, boards,
		`SELECT id, name, short_name
	FROM boards
	ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}

	return boards, nil
}

var ErrInvalidUserID = errors.New("invalid user id")
