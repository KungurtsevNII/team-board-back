package domain

import (
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
}

var (
	shortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

func NewBoard(name string, shortName string) (Board, error) {
	//TODO : доделать валиадцию
	if name == "" {
		return Board{}, InvalidNameErr
	}
	if len(name) > 100 {
		return Board{}, InvalidNameErr
	}

	if shortName == "" {
		return Board{}, InvalidNameErr
	}
	if !shortNameRegex.MatchString(shortName) {
		return Board{}, InvalidNameErr
	}

	now := time.Now().UTC()

	return Board{
		ID:        uuid.New(),
		Name:      name,
		ShortName: shortName,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}, nil
}

var InvalidNameErr = errors.New("invalid board name or short name")
