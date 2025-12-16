package puttask

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

type Repo interface {
	CheckColumnInBoard(ctx context.Context, boardID uuid.UUID, columnID uuid.UUID) (bool, error)
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
}

func (uc *UC) Handle(ctx context.Context, cmd Command) (task *domain.Task, err error) {
	foundDmn, err := uc.repo.GetTaskByID(ctx, cmd.TaskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	ex, err := uc.repo.CheckColumnInBoard(ctx, cmd.BoardID, cmd.ColumnID) 
	if err != nil {
		return nil, errors.Wrap(ErrPutTaskUnknown, err.Error())
	}
	if !ex {
		return nil, ErrColumnNotFound
	}

	foundDmn.Update(
		cmd.ColumnID,
		cmd.BoardID,
		cmd.Number,
		cmd.Title,
		cmd.Description,
		cmd.Tags,
		cmd.Checklists,
	)

	err = uc.repo.UpdateTask(ctx, foundDmn)
	if err != nil {
		return nil, errors.Wrap(ErrPutTaskUnknown, err.Error())
	}

	return foundDmn, nil
}