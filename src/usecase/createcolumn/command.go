package createcolumn

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateColumnCommand struct {
	BoardID  uuid.UUID `validate:"required,uuid"`
	Name     string `validate:"required,min=1,max=100"`
}

func NewCreateColumnCommand(boardID, name string) (CreateColumnCommand, error) {
	//TODO: Возможно стоит вынести как синглтон
	validate := validator.New()

	bID, err := uuid.Parse(boardID)
	if err != nil {
		return CreateColumnCommand{}, fmt.Errorf("%w: %v", ErrInvalidUUID, err)
	}

	ccc := CreateColumnCommand{
		BoardID:   bID,
		Name:      name,
	}

	err = validate.Struct(ccc)
	if err != nil {
		return CreateColumnCommand{}, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	return ccc, nil
}

type GetTaskStatusQuery struct {
	TaskID string
}
