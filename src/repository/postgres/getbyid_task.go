package postgres

import (
	"context"
	"encoding/json"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*domain.Task, error) {
    const op = "postgres.GetLastNumberTask"

    ds := goqu.From("tasks").
		Where(
			goqu.C("id").Eq(taskID),
			goqu.C("deleted_at").IsNull(),
		)
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var task TaskRecord
	err = pgxscan.Get(ctx, r.pool, &task, sql, params...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var cl []domain.Checklist
	err = json.Unmarshal(task.Checklists, &cl)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	dmn := domain.Task{
		ID : task.ID,
		ColumnID: task.ColumnID,
		BoardID: task.BoardID,
		Number: task.Number,
		Title: task.Title,
		Description: task.Description,
		Tags: task.Tags,
		Checklists: cl,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		DeletedAt: task.DeletedAt,
	}

    return &dmn, nil
}

