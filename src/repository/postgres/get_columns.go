package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetColumns(ctx context.Context, ID uuid.UUID) ([]domain.Column, error) {
	const op = "postgres.GetBoard"

	columns := make([]domain.Column, 0)
	err := pgxscan.Select(ctx, r.pool, &columns,
		`SELECT id, board_id, order_num, name 
		FROM columns WHERE board_id = $1 
		AND deleted_at IS NULL
		ORDER BY order_num;`, ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return columns, nil
}
