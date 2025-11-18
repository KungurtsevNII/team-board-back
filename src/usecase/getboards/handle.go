package getboards

import (
	"context"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repo interface {
	GetBoards(user_id uuid.UUID, ctx context.Context) ([]domain.Board, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(cmd GetBoardsQuery, ctx context.Context) ([]domain.Board, error) {
	const op = "getboards.Handle"

	boards, err := uc.repo.GetBoards(cmd.UserID, ctx)
	if err != nil {
		errors.Wrap(err, op)
	}

	return boards, nil
}
