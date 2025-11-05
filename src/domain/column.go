package domain

import (
	"time"
)

type Column struct {
	ID        string
	BoardID   string
	Name      string
	CreatedAt time.Time
	DeletedAt *time.Time
	UpdatedAt time.Time
}

func NewColumn(boardID string, name string) (*Column, error) {
	// todo validation

	return &Column{
		ID:        "adsasdas",
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
