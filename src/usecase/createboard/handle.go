package createboard

import (
	"context"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

type Repo interface {
	CheckBoard(shortName string, ctx context.Context) bool
	CreateBoard(board domain.Board, ctx context.Context) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(cmd CreateBoardCommand, ctx context.Context) (string, error) {
	const op = "createboard.Handle"
	if uc.repo.CheckBoard(cmd.ShortName, ctx) {
		return "", BoardIsExistsErr
	}

	board, err := domain.NewBoard(cmd.Name, cmd.ShortName)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	err = uc.repo.CreateBoard(board, ctx)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return board.ID.String(), nil
}
