package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	nameOfFirstColumn = "TODO"
)

var (
	ErrInvalidName = errors.New("invalid board name or short name")
	ErrColumnsIsEmpty = errors.New("columns is empty")
	shortNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,10}$`)
)

type Board struct {
	ID          uuid.UUID
	Name        string
	ShortName   string
	CreatedAt   time.Time
	DeletedAt   *time.Time
	UpdatedAt   time.Time
	Columns     []Column
	Tasks       []Task
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
	brdID := uuid.New()

	columns := make([]Column, 0)
	column, err := NewColumn(brdID, nameOfFirstColumn, 0)
	if err != nil {
		return Board{}, errors.Wrap(err, op)
	}
	columns = append(columns, *column)

	return Board{
		ID:        brdID,
		Name:      name,
		ShortName: shortName,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
		Columns: columns,
	}, nil
}

func (b *Board) GetFirstColumn() (Column, error) {
	if len(b.Columns) == 0 {
		return Column{}, ErrColumnsIsEmpty
	}
	return b.Columns[0], nil
}

func (c *Board) Delete() {
	now := time.Now().UTC()
	c.DeletedAt = &now
}
