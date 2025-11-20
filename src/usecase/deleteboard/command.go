package deleteboard

import (
	"github.com/google/uuid"
)

type Command struct {
	ID uuid.UUID
}

func NewCommand(id uuid.UUID) (Command, error) {
	if id == uuid.Nil {
		return Command{}, ErrBoardIdEmpty
	}

	_, err := uuid.Parse(id.String())
	if err != nil {
		return Command{}, ErrBoardIdInvalid
	}

	return Command{
		ID: id,
	}, nil
}
