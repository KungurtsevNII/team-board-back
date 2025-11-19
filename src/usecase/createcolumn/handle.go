package createcolumn

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	// "github.com/KungurtsevNII/team-board-back/src/repository/postgres"
)

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

type Repo interface {
	CheckBoard(id string, ctx context.Context) bool
	GetLastOrderNumColumn(
		ctx context.Context,
		boardID uuid.UUID,
	) (orderNum int64, err error)
	CreateColumn(
		ctx context.Context,
		column *domain.Column,
	) (err error)
}

func (uc *UC) Handle(ctx context.Context, cmd CreateColumnCommand) (column *domain.Column, err error) {
	// todo errrors wrap
	if !uc.repo.CheckBoard(cmd.BoardID.String(), ctx) {
		return nil, fmt.Errorf("%w: %v", ErrBoardIsNotExists, err)
	}

	orderNum, err := uc.repo.GetLastOrderNumColumn(ctx, cmd.BoardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			orderNum = 0
		} else {
			// todo errors wrap
			return nil, fmt.Errorf("%w: %v", ErrGetLastOrderNumUnknown, err)
		}
	} else {
		orderNum++
	}

	column, err = domain.NewColumn(cmd.BoardID, cmd.Name, orderNum)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidationFailed, err)
	}

	err = uc.repo.CreateColumn(ctx, column)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateColumnUnknown, err)
	}

	return column, nil
}
