package postgres

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) CheckColumnInBoard(ctx context.Context, boardID uuid.UUID, columnID uuid.UUID) (bool, error) {
	const op = "postgres.CheckColumnInBoard"
	
	ds := goqu.From("columns").
		Select(goqu.COUNT("*")).
		Where(
			goqu.C("board_id").Eq(boardID), 
			goqu.C("id").Eq(columnID),
		)
	
	sql, params, err := ds.ToSQL()
	if err != nil {
        return false, errors.Wrap(err, op)
	}
	
	var count int64
	err = pgxscan.Get(ctx, r.pool, &count, sql, params...)
	if err != nil {
        return false, errors.Wrap(err, op)
	}
	
	return count > 0, nil
}
