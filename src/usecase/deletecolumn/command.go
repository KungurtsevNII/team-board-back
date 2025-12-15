package deletecolumn

import (
	"github.com/google/uuid"
)

type Command struct {
	ColumnID uuid.UUID 
}

func NewCommand(taskID string) (Command, error) {
	uid, err := uuid.Parse(taskID)
	if err != nil {
		return Command{}, ErrInvalidColumnID
	}

	return Command{
		ColumnID: uid,
	}, nil
}
