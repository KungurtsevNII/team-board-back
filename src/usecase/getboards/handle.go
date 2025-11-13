package getboards

import (
	"context"
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
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

func (uc *UC) Handle(cmd GetBoardsCommand, ctx context.Context) ([]domain.Board, error) {
	const op = "getboards.Handle"

	boards, err := uc.repo.GetBoards(cmd.UserID, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return boards, nil
}
