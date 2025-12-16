package puttask

import (
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Command struct {
	TaskID      uuid.UUID `validate:"required,uuid"`
	BoardID     uuid.UUID `validate:"required,uuid"`
	ColumnID    uuid.UUID `validate:"required,uuid"`
	Title       string    `validate:"required,min=1,max=255"`
	Number      int64     
	Description *string
	Tags        []string
	Checklists  []domain.Checklist
}

func NewCommand(
	taskID, boardID, columnID, title string,
	number int64, description *string,
	tags []string, checklists []domain.Checklist,
) (Command, error) {
	validate := validator.New()

	tID, err := uuid.Parse(taskID)
	if err != nil {
		return Command{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	bID, err := uuid.Parse(boardID)
	if err != nil {
		return Command{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	cID, err := uuid.Parse(columnID)
	if err != nil {
		return Command{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	cmd := Command{
		TaskID:      tID,
		BoardID:     bID,
		ColumnID:    cID,
		Number:      number,
		Title:       title,
		Description: description,
		Tags:        tags,
		Checklists:  checklists,
	}

	err = validate.Struct(cmd)
	if err != nil {
		return Command{}, errors.Wrap(ErrValidationFailed, err.Error())
	}

	return cmd, nil
}
