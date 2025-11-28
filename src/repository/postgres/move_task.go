package postgres

import (
	"context"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) MoveTaskColumn(ctx context.Context, taskID uuid.UUID, columnID uuid.UUID) (*domain.Task, error) {
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
		).
		Returning(goqu.Star())

	sql, params, err := ds.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var task TaskRecord
	err = pgxscan.Get(ctx, r.pool, &task, sql, params...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	dmn, err := task.toDomain()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return dmn, nil
}
