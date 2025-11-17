package domain

import (
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID
	BoardID   string
	Name      string
	OrderNum  int
	CreatedAt time.Time
	DeletedAt *time.Time
	UpdatedAt time.Time
	Tasks     []any //TODO
}

func NewColumn(boardID string, name string) (*Column, error) {
	// todo validation

	return &Column{
		ID:        uuid.New(),
		BoardID:   boardID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}, nil
}

func (c *Column) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}
