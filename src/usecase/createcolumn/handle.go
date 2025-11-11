package createcolumn

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	CheckBoard(id string) bool 
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
	const op = "createcolumn.Handle"
	log := slog.Default().With("op", op, "cmd", cmd)

	if !uc.repo.CheckBoard(cmd.BoardID.String()){
	    return nil, fmt.Errorf("%s: %v", op, ErrBoardIsNotExistsErr)
	}

	orderNum, err := uc.repo.GetLastOrderNumColumn(ctx, cmd.BoardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			orderNum = 0
		}else{
			log.Error("failed to get last order num", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %v", op, ErrGetLastOrderNumErr)
		}
	}else{
		orderNum++
	}

	column, err = domain.NewColumn(cmd.BoardID, cmd.Name, orderNum)
	if err != nil {
		log.Error("failed to create domain column", slog.String("err", err.Error()))
	    return nil, fmt.Errorf("%s: %v", op, ErrValidationFailed)
	}

	err = uc.repo.CreateColumn(ctx, column)
	if err != nil {
		log.Error("failed to create column", slog.String("err", err.Error()))
	    return nil, fmt.Errorf("%s: %v", op, ErrCreateColumnErr)
	}

	return column, nil
}