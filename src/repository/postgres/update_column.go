package postgres

import (
	"context"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

func (r Repository) UpdateColumn(ctx context.Context, column *domain.Column) error {
	const op = "postgres.UpdateColumn"

	column.UpdatedAt = time.Now().UTC()


	ds := goqu.Update("columns").Where(
		goqu.C("id").Eq(column.ID),
		goqu.C("deleted_at").IsNull(),
	).Set(
		ColumnRecord{
			BoardID: column.BoardID,
			Name:        column.Name,
			OrderNum: column.OrderNum,
			UpdatedAt:   column.UpdatedAt,
			DeletedAt: column.DeletedAt,
		},
	)
	sql, params, err := ds.ToSQL()
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = r.pool.Exec(ctx, sql, params...)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
