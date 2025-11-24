package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrInvalidName = errors.New("invalid board name or short name")
	shortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

type Board struct {
	ID        uuid.UUID
	Name      string
	ShortName string
	CreatedAt time.Time
	DeletedAt *time.Time
	UpdatedAt time.Time
}

func NewBoard(name string, shortName string) (Board, error) {
	const op = "domain.NewBoard"
	if name == "" {
		return Board{}, errors.Wrap(ErrInvalidName, op)
	}
	if len(name) > 100 {
		return Board{}, errors.Wrap(ErrInvalidName, op)
	}

	if shortName == "" {
		return Board{}, errors.Wrap(ErrInvalidName, op)
	}
	if !shortNameRegex.MatchString(shortName) {
		return Board{}, errors.Wrap(ErrInvalidName, op)
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
