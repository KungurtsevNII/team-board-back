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

	targetDmn := &domain.Task{
		ID:          cmd.TaskID,
		ColumnID:    cmd.ColumnID,
		BoardID:     cmd.BoardID,
		Number:      cmd.Number,
		Title:       cmd.Title,
		Description: cmd.Description,
		Tags:        cmd.Tags,
		Checklists:  cmd.Checklists,
		CreatedAt:   foundDmn.CreatedAt,
		UpdatedAt:   foundDmn.UpdatedAt,
	}

	err = uc.repo.UpdateTask(ctx, targetDmn)
	if err != nil {
		return nil, errors.Wrap(ErrPutTaskUnknown, err.Error())
	}

	return targetDmn, nil
}