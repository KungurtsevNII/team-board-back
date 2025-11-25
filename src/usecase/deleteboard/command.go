package deleteboard

import (
	"github.com/google/uuid"
)

type Command struct {
	ID uuid.UUID
}

func NewCommand(id string) (Command, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return Command{}, ErrBoardIdInvalid
	}

	return Command{
		ID: uid,
	}, nil
}
