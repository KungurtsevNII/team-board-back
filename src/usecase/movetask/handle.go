package movetask

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

func (uc *UC) Handle(ctx context.Context, cmd MoveTaskCommand) (*domain.Task, error) {
	task, err := uc.repo.GetTaskByID(ctx, cmd.TaskID)
	if err != nil {
		return nil, errors.Wrap(ErrTaskNotFound, err.Error())
	}

	ex, err := uc.repo.CheckColumnInBoard(ctx, task.BoardID, cmd.ColumnID)
	if err != nil {
		return nil, errors.Wrap(ErrMoveTaskUnknown, err.Error())
	}
	if !ex {
		return nil, ErrColumnNotInBoard
	}

	err = task.MoveToColumn(cmd.ColumnID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyInColumn) {
			return task, nil
		}
		return nil, errors.Wrap(ErrMoveTaskUnknown, err.Error())
	}

	err = uc.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, errors.Wrap(ErrMoveTaskUnknown, err.Error())
	}

	return task, nil
}
