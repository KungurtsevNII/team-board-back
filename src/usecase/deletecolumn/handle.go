package deletecolumn

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Repo interface {
	GetColumnByID(ctx context.Context, columnID uuid.UUID) (*domain.Column, error)
	CheckColumnIsEmpty(ctx context.Context, columnID uuid.UUID) (bool, error)
	UpdateColumn(ctx context.Context, column *domain.Column) error
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) error{
	dmn, err := uc.repo.GetColumnByID(ctx, cmd.ColumnID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrColumnNotFound
		}
		return errors.Wrap(ErrGetColumnUnknown, err.Error())
	}

	isEmpty, err := uc.repo.CheckColumnIsEmpty(ctx, cmd.ColumnID)
	if err != nil {
		return errors.Wrap(ErrCheckColumnIsEmptyUnknown, err.Error())
	}
	if !isEmpty {
		return ErrColumnNotEmpty
	}

	dmn.Delete()

	err = uc.repo.UpdateColumn(ctx, dmn)
	if err != nil {
		return errors.Wrap(ErrDeleteColumnUnknown, err.Error())
	}
	return nil
}
