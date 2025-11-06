package createcolumn

import (
	"fmt"

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
	CheckColumn(name string) bool
	CreateColumn(column domain.Column) error
}

func (uc *UC) Handle(cmd CreateColumnCommand) error {
	//тут прям обращение к базе данных
	const op = "createcolumn.Handle"

	if uc.repo.CheckColumn(cmd.Title) {
	    return ColumnIsExistsErr
	}

	column, err := domain.NewColumn(cmd.Title, cmd.BoardID)
	if err != nil {
	    return fmt.Errorf("%s: %v", op, err)
	}
	return uc.repo.CreateColumn(*column)
}