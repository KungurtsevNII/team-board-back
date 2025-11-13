package domain

import (
	// "errors"
	// "strings"
	"errors"
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
	Columns   []Column
}

var (
	shortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewBoard(name string, shortName string) (*Board, error) {
	//todo validation
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

var ErrBoardNotFound = errors.New("board not found")
var ErrInvalidID = errors.New("invalid id format")
