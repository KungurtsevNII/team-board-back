package gettask

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
}

func (uc *UC) Handle(ctx context.Context, query GetTaskQuery) (*domain.Task, error) {
	dmn, err := uc.repo.GetTaskByID(ctx, query.TaskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, errors.Wrap(ErrGetTaskUnknown, err.Error())
	}

	return dmn, nil
}