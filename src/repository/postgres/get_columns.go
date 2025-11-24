package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetColumns(ctx context.Context, ID string) ([]domain.Column, error) {
	const op = "postgres.GetBoard"
	uid, err := uuid.Parse(ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	columns := make([]domain.Column, 0)
	err = pgxscan.Select(ctx, r.pool, columns,
		`SELECT id, board_id, order_num, name 
		FROM columns WHERE board_id = $1 
		ORDER BY order_num;`, uid)
	if err != nil {
		return nil, err
	}
	return columns, nil
}
