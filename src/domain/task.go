package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrAlreadyInColumn = errors.New("task already in target column")

type Task struct {
	ID          uuid.UUID
	ColumnID    uuid.UUID
	ColumnName  *string
	BoardID     uuid.UUID
	BoardName   *string
	BoardShortName *string
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
		DeletedAt:   nil,
	}, nil
}

func (t *Task)Update(
	columnID uuid.UUID,
	boardID uuid.UUID,
	number int64,
	title string,
	description *string,
	tags []string,
	checklists []Checklist,
){
	t.ColumnID = columnID
	t.BoardID = boardID
	t.Number = number
	t.Title = title
	t.Description = description
	t.Tags = tags
	t.Checklists = checklists
	t.UpdatedAt = time.Now().UTC()
}

func (c *Task) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}

func (t *Task) MoveToColumn(columnID uuid.UUID) error {
	if t.ColumnID == columnID {
		return ErrAlreadyInColumn
	}

	t.ColumnID = columnID
	t.UpdatedAt = time.Now().UTC()
	return nil
}
