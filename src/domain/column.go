package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrEmptyColumnName = errors.New("column name can't be empty")

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
	if name == "" {
		return nil, ErrEmptyColumnName
	}
	id := uuid.New()

	return &Column{
		ID:        id,
		BoardID:   boardID,
		Name:      name,
		OrderNum:  orderNum,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
	}, nil
}

func (c *Column) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}
