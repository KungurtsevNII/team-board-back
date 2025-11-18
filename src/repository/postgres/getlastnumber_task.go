package postgres

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetLastNumberTask(ctx context.Context, boardID uuid.UUID) (int64, error) {
    const op = "postgres.GetLastNumberTask"

    ds := goqu.From("tasks").
		Where(goqu.C("board_id").Eq(boardID)).
		Select(goqu.C("number")).
		Order(goqu.I("number").Desc()).Limit(1)
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	var number int64
	err = pgxscan.Get(ctx, r.pool, &number, sql, params...)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

    return number, nil
}

