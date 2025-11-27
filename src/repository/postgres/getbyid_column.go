package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetColumnByID(ctx context.Context, columnID uuid.UUID) (*domain.Column, error) {
    const op = "postgres.GetColumnByID"

    ds := goqu.From("columns").
		Where(
			goqu.C("id").Eq(columnID),
			goqu.C("deleted_at").IsNull(),
		)
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var column ColumnRecord
	err = pgxscan.Get(ctx, r.pool, &column, sql, params...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	dmn, err := column.toDomain()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

    return dmn, nil
}

