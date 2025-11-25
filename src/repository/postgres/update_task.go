package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r Repository) UpdateTask(ctx context.Context, task *domain.Task) error {
	const op = "postgres.UpdateTask"

	task.UpdatedAt = time.Now().UTC()

	checklistsJSON, err := json.Marshal(task.Checklists)
	if err != nil {
		return errors.Wrap(err, op)
	}

	var tagsValue interface{}
	if task.Tags != nil {
		tagsValue = pq.StringArray(task.Tags)
	} else {
		tagsValue = nil
	}

	ds := goqu.Update("tasks").Where(
		goqu.C("id").Eq(task.ID),
		goqu.C("deleted_at").IsNull(),
	).Set(
		goqu.Record{
			"board_id":    task.BoardID,
			"column_id":   task.ColumnID,
			"number":      task.Number,
			"title":       task.Title,
			"description": task.Description,
			"tags":        tagsValue,
			"checklists":  checklistsJSON,
			"updated_at":  task.UpdatedAt,
			"deleted_at":  task.DeletedAt,
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
