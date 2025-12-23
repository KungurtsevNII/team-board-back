package createcolumn

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

//go generate 
type Repo interface {
	CheckBoard(ctx context.Context, id string) bool
	GetLastOrderNumColumn(
		ctx context.Context,
		boardID uuid.UUID,
	) (orderNum int64, err error)
	CreateColumn(
		ctx context.Context,
		column *domain.Column,
	) (err error)
}

func (uc *UC) Handle(ctx context.Context, cmd Command) (column *domain.Column, err error) {
	if !uc.repo.CheckBoard(ctx, cmd.BoardID.String()) {
		return nil, ErrBoardIsNotExists
	}

	orderNum, err := uc.repo.GetLastOrderNumColumn(ctx, cmd.BoardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			orderNum = 0
		} else {
			return nil, errors.Wrap(ErrGetLastOrderNumUnknown, err.Error())
		}
	} else {
		orderNum++
	}

	column, err = domain.NewColumn(cmd.BoardID, cmd.Name, orderNum)
	if err != nil {
		return nil, errors.Wrap(ErrValidationFailed, err.Error())
	}

	err = uc.repo.CreateColumn(ctx, column)
	if err != nil {
		return nil, errors.Wrap(ErrCreateColumnUnknown, err.Error())
	}

	return column, nil
}
