package postgres

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) DeleteTask(ctx context.Context, taskID uuid.UUID) error{
    const op = "postgres.DeleteTask"
	
	nowUTC := time.Now().UTC()
    ds := goqu.Update("tasks").
        Where(goqu.C("id").Eq(taskID)).
        Set(goqu.Record{
            "deleted_at": nowUTC,
        })
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = r.pool.Exec(ctx, sql, params...)
    return err
}

