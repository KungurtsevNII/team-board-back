package postgres

import (
	"context"
	"encoding/json"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

func (r Repository) CreateTask(ctx context.Context, task *domain.Task) error {
	const op = "postgres.CreateTask"
	
	sql := `INSERT INTO tasks (
		id, board_id, column_id, number, title, description, tags,
		checklists, created_at, updated_at, deleted_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	checklistsJSON, err := json.Marshal(task.Checklists)
	if err != nil {
		return errors.Wrap(err, op)
	}

	taskRecord := TaskRecord{
		ID:          task.ID,
		BoardID:     task.BoardID,
		ColumnID:    task.ColumnID,
		Number:      task.Number,
		Title:       task.Title,
		Description: task.Description,
		Tags:        task.Tags,
		Checklists:  checklistsJSON,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		DeletedAt:   task.DeletedAt,
	}

	_, err = r.pool.Exec(ctx, sql,
		taskRecord.ID,
		taskRecord.BoardID,
		taskRecord.ColumnID,
		taskRecord.Number,
		taskRecord.Title,
		taskRecord.Description,
		taskRecord.Tags,
		taskRecord.Checklists,
		taskRecord.CreatedAt,
		taskRecord.UpdatedAt,
		taskRecord.DeletedAt,
	)

	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
