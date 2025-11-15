package createtask

import (
	"context"


	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	// "github.com/KungurtsevNII/team-board-back/src/repository/postgres"
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


type Repo interface {
	CheckColumnInBoard(ctx context.Context, boardID uuid.UUID, columnID uuid.UUID) (bool,error)
	GetLastNumberTask(ctx context.Context, boardID uuid.UUID) (int64, error)
	CreateTask(ctx context.Context, task *domain.Task) error
}

func (uc *UC) Handle(ctx context.Context, cmd CreateTaskCommand) (task *domain.Task, err error) {
	ex, err := uc.repo.CheckColumnInBoard(ctx, cmd.BoardID, cmd.ColumnID) 
	if err != nil {
		return nil, errors.Wrap(ErrCheckColumnInBoardFailed, err.Error())
	}
	if !ex{
		return nil, ErrColumnOrBoardIsNotExists
	}

	number, err := uc.repo.GetLastNumberTask(ctx, cmd.BoardID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			number = 0
		}else{
			return nil, errors.Wrap(ErrGetLastNumberFailed, err.Error())
		}
	}else{
		number++
	}

	task, err = domain.NewTask(
		cmd.ColumnID,
		cmd.BoardID,
		number,
		cmd.Title,
		cmd.Description,
		cmd.Tags,
		cmd.Checklists,
	)
	if err != nil {
		return nil, errors.Wrap(ErrValidationFailed, err.Error())
	}

	err = uc.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, errors.Wrap(ErrCreateTaskUnknown, err.Error())
	}

	return task, nil
}