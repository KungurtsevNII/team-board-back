package domain

import (
	"time"

	"github.com/google/uuid"
	"encoding/json"
)

type Task struct {
	ID          uuid.UUID
	ColumnID    uuid.UUID
	BoardID     uuid.UUID
	Number      int64
	Title       string
	Description *string
	Tags        []string
	Checklists  *json.RawMessage
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
	checklists *json.RawMessage,
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

func (c *Task) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}
