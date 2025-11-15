package createtask

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateTaskCommand struct {
	ColumnID    uuid.UUID `validate:"required,uuid"`
	BoardID     uuid.UUID `validate:"required,uuid"`
	Title       string    `validate:"required,min=1,max=255"`
	Description *string
	Tags        []string
	Checklists  *json.RawMessage
}

func NewCreateTaskCommand(
	columnID, boardID, name string,
	description *string,
	tags []string,
	checklists *json.RawMessage,
) (CreateTaskCommand, error) {

	validate := validator.New()

	bID, err := uuid.Parse(boardID)
	if err != nil {
		return CreateTaskCommand{}, errors.Wrap(ErrInvalidUUID, err.Error())
	}

	cID, err := uuid.Parse(columnID)
	if err != nil {
		return CreateTaskCommand{}, errors.Wrap(ErrInvalidUUID, err.Error())
	}

	ctc := CreateTaskCommand{
		ColumnID:    cID,
		BoardID:     bID,
		Title:       name,
		Description: description,
		Tags:        tags,
		Checklists:  checklists,
	}

	err = validate.Struct(ctc)
	if err != nil {
		return CreateTaskCommand{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	return ctc, nil
}

type GetTaskStatusQuery struct {
	TaskID string
}
