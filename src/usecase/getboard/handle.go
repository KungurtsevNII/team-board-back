package getboard

import (
	"context"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

type Repo interface {
	GetBoard(ctx context.Context, ID string) (domain.Board, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd GetBoardCommand) (domain.Board, error) {
	const op = "getboard.Handle"

	board, err := uc.repo.GetBoard(ctx, cmd.ID)
	if err != nil {
		return domain.Board{}, errors.Wrap(err, op)
	}
	return board, nil
}
