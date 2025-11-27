package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r Repository) CheckColumnIsEmpty(ctx context.Context, columnID uuid.UUID) (bool, error){
	const op = "postgres.CheckColumnIsEmpty"
	
	ds := goqu.From("tasks").
		Select(goqu.COUNT("*")).
		Where(
			goqu.C("column_id").Eq(columnID),
			goqu.C("deleted_at").IsNull(),
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
	
	return count == 0, nil
}