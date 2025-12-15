package deletetask

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) error{
	dmn, err := uc.repo.GetTaskByID(ctx, cmd.TaskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return errors.Wrap(ErrGetTaskUnknown, err.Error())
	}

	dmn.Delete()
	err = uc.repo.UpdateTask(ctx, dmn)
	if err != nil {
		return errors.Wrap(ErrDeleteTaskUnknown, err.Error())
	}
	return nil
}
