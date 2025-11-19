package postgres

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) MoveTaskColumn(ctx context.Context, taskID uuid.UUID, columnID uuid.UUID) error {
	const op = "postgres.MoveTaskColumn"

	now := time.Now().UTC()

	ds := goqu.Update("tasks").
		Set(goqu.Record{
			"column_id":  columnID,
			"updated_at": now,
		}).
		Where(
			goqu.C("id").Eq(taskID),
			goqu.C("deleted_at").IsNull(),
		)

	sql, params, err := ds.ToSQL()
	if err != nil {
		return errors.Wrap(err, op)
	}
	result, err := r.pool.Exec(ctx, sql, params...)
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Проверяем, что задача действительно была обновлена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTaskNotFoundOrDeleted
	}

	return nil
}
