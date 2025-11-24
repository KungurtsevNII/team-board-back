package getboards

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repo interface {
	GetBoards(ctx context.Context, user_id uuid.UUID) ([]domain.Board, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Query) ([]domain.Board, error) {
	boards, err := uc.repo.GetBoards(ctx, cmd.UserID)
	if err != nil {
		errors.Wrap(ErrGetBoards, err.Error())
	}

	return boards, nil
}
