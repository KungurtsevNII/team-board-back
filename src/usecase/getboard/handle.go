package getboard

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	GetBoard(ctx context.Context, ID uuid.UUID) (*domain.Board, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, quer Query) (*domain.Board, error) {
	const op = "getboard.Handle"

	board, err := uc.repo.GetBoard(ctx, quer.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, errors.Wrap(ErrBoardNotFound, err.Error())
		}
		return nil, errors.Wrap(err, op)
	}
	return board, nil
}
