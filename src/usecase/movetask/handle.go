package movetask

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
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
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*domain.Task, error)
	CheckColumnInBoard(ctx context.Context, boardID uuid.UUID, columnID uuid.UUID) (bool, error)
	MoveTaskColumn(ctx context.Context, taskID uuid.UUID, columnID uuid.UUID) error
}

func (uc *UC) Handle(ctx context.Context, cmd MoveTaskCommand) (*domain.Task, error) {
	task, err := uc.repo.GetTaskByID(ctx, cmd.TaskID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if task.ColumnID == cmd.ColumnID {
		return task, nil
	}

	// TODO: Проверять существование колонки

	// Проверяем, что целевая колонка принадлежит той же доске
	ok, err := uc.repo.CheckColumnInBoard(ctx, task.BoardID, cmd.ColumnID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check column in board")
	}
	if !ok {
		return nil, ErrColumnNotInBoard
	}

	// Обновляем колонку задачи
	if err := uc.repo.MoveTaskColumn(ctx, cmd.TaskID, cmd.ColumnID); err != nil {
		if errors.Is(err, postgres.ErrTaskNotFoundOrDeleted) {
			return nil, ErrTaskNotFound
		}
		return nil, ErrMoveTaskUnknown
	}

	task.ColumnID = cmd.ColumnID
	return task, nil
}
