package deleteboard

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repo interface {
	DeleteBoard(ctx context.Context, id uuid.UUID) error
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
	const op = "deleteboard.Handle"

	err := uc.repo.DeleteBoard(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
