package domain

import (
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID
	BoardID   uuid.UUID
	OrderNum  int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewColumn(boardID uuid.UUID, name string, orderNum int64) (*Column, error) {
	id := uuid.New()

	return &Column{
		ID:        id,
		BoardID:   boardID,
		Name:      name,
		OrderNum:  orderNum,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}, nil
}

func (c *Column) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}
