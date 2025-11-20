package createboard

import (
	"context"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

type Repo interface {
	CheckBoard(ctx context.Context, shortName string) bool
	CreateBoard(ctx context.Context, board domain.Board) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) (string, error) {
	if uc.repo.CheckBoard(ctx, cmd.ShortName) {
		return "", ErrBoardIsExists
	}

	board, err := domain.NewBoard(cmd.Name, cmd.ShortName)
	if err != nil {
		return "", errors.Wrap(ErrNewBoardFailed, err.Error())
	}

	err = uc.repo.CreateBoard(ctx, board)
	if err != nil {
		return "", errors.Wrap(ErrCreateBoard, err.Error())
	}

	return board.ID.String(), nil
}
