package createboard

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Repo interface {
	CheckBoard(ctx context.Context, shortName string) bool
	CreateBoard(ctx context.Context, board domain.Board) error
	GetLastOrderNumColumn(
		ctx context.Context,
		boardID uuid.UUID,
	) (orderNum int64, err error)
	CreateColumn(
		ctx context.Context,
		column *domain.Column,
	) (err error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, cmd Command) (*domain.Board, *domain.Column, error) {
	if uc.repo.CheckBoard(ctx, cmd.ShortName) {
		return nil, nil, ErrBoardIsExists
	}

	board, err := domain.NewBoard(cmd.Name, cmd.ShortName)
	if err != nil {
		return nil, nil, errors.Wrap(ErrValidationFailed, err.Error())
	}

	number, err := uc.repo.GetLastOrderNumColumn(ctx, board.ID)
	number++
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, errors.Wrap(ErrGetLastOrderNumUnknown, err.Error())
	}
	if errors.Is(err, pgx.ErrNoRows) {
		number = 0
	}

	column, err := domain.NewColumn(board.ID, cmd.ColumnName, number)
	if err != nil {
		return nil, nil, errors.Wrap(ErrValidationFailed, err.Error())
	}

	err = uc.repo.CreateBoard(ctx, board)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCreateBoard, err.Error())
	}

	err = uc.repo.CreateColumn(ctx, column)
	if err != nil {
		return nil, nil, errors.Wrap(ErrCreateColumnUnknown, err.Error())
	}

	return &board, column, nil
}
