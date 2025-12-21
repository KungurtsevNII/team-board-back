package createcolumn

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Command struct {
	BoardID uuid.UUID `validate:"required,uuid"`
	Name    string    `validate:"required,min=1,max=100"`
}

func NewCommand(boardID, name string) (Command, error) {
	validate := validator.New()

	bID, err := uuid.Parse(boardID)
	if err != nil {
		return Command{}, errors.Wrap(ErrInvalidUUID, err.Error())
	}

	ccc := Command{
		BoardID: bID,
		Name:    name,
	}

	err = validate.Struct(ccc)
	if err != nil {
		return Command{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	return ccc, nil
}
