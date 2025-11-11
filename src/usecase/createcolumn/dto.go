package createcolumn

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateColumnCommand struct {
	BoardID  uuid.UUID `validate:"required,uuid"`
	Name     string `validate:"required,min=1,max=100"`
}

func NewCreateColumnCommand(boardID, name string) (CreateColumnCommand, error) {
	const op = "createcolumn.NewCreateColumnCommand"
	log := slog.Default().With("op", op,"boardID", boardID, "name", name,)
	//TODO: Возможно стоит вынести как синглтон
	validate := validator.New()

	bID, err := uuid.Parse(boardID)
	if err != nil {
		log.Warn("failed to parse boardID", "err", err)
		return CreateColumnCommand{}, fmt.Errorf("%s: %w", op, ErrInvalidUUID)
	}

	ccc := CreateColumnCommand{
		BoardID:   bID,
		Name:      name,
	}

	err = validate.Struct(ccc)
	if err != nil {
		log.Warn("validation failed", "err", err)
		return CreateColumnCommand{}, fmt.Errorf("%s: %w", op, ErrValidationFailed)
	}

	return ccc, nil
}

type GetTaskStatusQuery struct {
	TaskID string
}
