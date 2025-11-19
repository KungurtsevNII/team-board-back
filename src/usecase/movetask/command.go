package movetask

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MoveTaskCommand struct {
	TaskID   uuid.UUID `validate:"required,uuid"`
	ColumnID uuid.UUID `validate:"required,uuid"`
}

func NewMoveTaskCommand(taskID, newColumnID string) (MoveTaskCommand, error) {
	validate := validator.New()

	tID, err := uuid.Parse(taskID)
	if err != nil {
		return MoveTaskCommand{}, errors.Wrap(ErrInvalidUUID, err.Error())
	}

	cID, err := uuid.Parse(newColumnID)
	if err != nil {
		return MoveTaskCommand{}, errors.Wrap(ErrInvalidUUID, err.Error())
	}

	mtc := MoveTaskCommand{
		TaskID:   tID,
		ColumnID: cID,
	}

	err = validate.Struct(mtc)
	if err != nil {
		return MoveTaskCommand{}, errors.Wrap(ErrValidationFailed, err.Error())
	}
	return mtc, nil
}
