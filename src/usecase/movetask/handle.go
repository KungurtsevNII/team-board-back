package movetask

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	MoveTaskColumn(ctx context.Context, taskID uuid.UUID, columnID uuid.UUID) (*domain.Task, error)
}

func (uc *UC) Handle(ctx context.Context, cmd MoveTaskCommand) (*domain.Task, error) {
	task, err := uc.repo.GetTaskByID(ctx, cmd.TaskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, errors.Wrap(err, "failed to get task")
	}

	if task.ColumnID == cmd.ColumnID {
		return task, nil
	}

	// TODO: Проверять существование колонки

	// Получаем актуальное состояние после UPDATE
	updatedTask, err := uc.repo.MoveTaskColumn(ctx, cmd.TaskID, cmd.ColumnID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, errors.Wrap(ErrMoveTaskUnknown, err.Error())
	}

	return updatedTask, nil
}
