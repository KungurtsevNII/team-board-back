package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	ColumnID    uuid.UUID
	BoardID     uuid.UUID
	Number      int64
	Title       string
	Description *string
	Tags        []string
	Checklists  []Checklist
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewTask(
	columnID uuid.UUID,
	boardID uuid.UUID,
	number int64,
	title string,
	description *string,
	tags []string,
	checklists []Checklist,
) (*Task, error) {
	id := uuid.New()

	return &Task{
		ID:          id,
		ColumnID:    columnID,
		BoardID:     boardID,
		Number:      number,
		Title:       title,
		Description: description,
		Tags:        tags,
		Checklists:  checklists,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		DeletedAt:    nil,
	}, nil
}