package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetBoards(ctx context.Context, user_id uuid.UUID) ([]domain.Board, error) {
	const op = "postgres.GetBoards"
	//TODO проверять user_id и подставлять его как параметр в запрос
	_, err := uuid.Parse(user_id.String())
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	boards := make([]domain.Board, 0)
	//TODO подставлять user_id в запрос
	err = pgxscan.Select(ctx, r.pool, boards,
		`SELECT id, name, short_name
	FROM boards
	ORDER BY updated_at DESC`)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return boards, nil
}
