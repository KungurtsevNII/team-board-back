package deleteboard

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/jackc/pgx/v5"
)

type Repo interface {
	GetBoard(ctx context.Context, ID uuid.UUID) (*domain.Board, error)
	UpdateBoard(ctx context.Context, board *domain.Board) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) error {
	dmn, err := uc.repo.GetBoard(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrBoardDoesntExist
		}
		return errors.Wrap(ErrBoardDoesntExist, err.Error())
	}

	dmn.Delete()

	err = uc.repo.UpdateBoard(ctx, dmn)
	if err != nil {
		return errors.Wrap(ErrBoardDeleteUnknown, err.Error())
	}
	return nil
}
