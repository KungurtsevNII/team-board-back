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

func (uc *UC) Handle(cmd CreateBoardCommand, ctx context.Context) (*domain.Board, error) {
	if uc.repo.CheckBoard(cmd.ShortName, ctx) {
		return nil, ErrBoardIsExists
	}

	board, err := domain.NewBoard(cmd.Name, cmd.ShortName)
	if err != nil {
		return nil, errors.Wrap(ErrValidationFaild, err.Error())
	}

	err = uc.repo.CreateBoard(board, ctx)
	if err != nil {
		return nil, errors.Wrap(ErrCreateBoardFailed, err.Error())
	}

	return &board, nil
}
