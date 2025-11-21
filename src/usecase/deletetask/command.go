package deletetask

import (
	"github.com/google/uuid"
)

type Command struct {
	TaskID uuid.UUID 
}

func NewCommand(taskID string) (Command, error) {
	uid, err := uuid.Parse(taskID)
	if err != nil {
		return Command{}, ErrInvalidTaskID
	}

	return Command{
		TaskID: uid,
	}, nil
}
