package domain

import (
	// "errors"
	// "strings"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	CreatedAt time.Time
	DeletedAt *time.Time
	UpdatedAt time.Time
}

var (
	shortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewBoard(name string, shortName string) (*Board, error) {
	//тут не вижу смысла делать валидацию , т.к. она уже есть в dto
	now := time.Now()

	return &Board{
		ID:        uuid.New(),
		Name:      name,
		ShortName: shortName,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}, nil
}
