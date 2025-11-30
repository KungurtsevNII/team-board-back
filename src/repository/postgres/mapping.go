package postgres

import (
	"encoding/json"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

func (task *TaskRecord) toDomain() (*domain.Task, error) {
	const op = "postgres.TaskRecord.ToDomain"

	var cl []domain.Checklist
	err := json.Unmarshal(task.Checklists, &cl)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	dmn := domain.Task{
		ID:          task.ID,
		ColumnID:    task.ColumnID,
		BoardID:     task.BoardID,
		Number:      task.Number,
		Title:       task.Title,
		Description: task.Description,
		Tags:        task.Tags,
		Checklists:  cl,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		DeletedAt:   task.DeletedAt,
	}
	return &dmn, nil
}

func (col *ColumnRecord) toDomain() (*domain.Column, error) {
	return &domain.Column{
		ID:        col.ID,
		BoardID:   col.BoardID,
		OrderNum:  col.OrderNum,
		Name:      col.Name,
		CreatedAt: col.CreatedAt,
		UpdatedAt: col.UpdatedAt,
		DeletedAt: col.DeletedAt,
	}, nil
}
