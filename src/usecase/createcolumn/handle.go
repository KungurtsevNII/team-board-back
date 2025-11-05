package createcolumn

import (
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
)

type UC struct {
	repo postgres.Repository
}

type Repo interface {
	// Check()
	CreateColumn(column domain.Column) error
}

func (uc *UC) CreateColumnHandle(cmd CreateColumnCommand) error {
	//тут прям обращение к базе данных
	const op = "createcolumn.Handle"
	column, err := domain.NewColumn(cmd.Title, cmd.BoardID)
	if err != nil {
	    return fmt.Errorf("%s: %v", op, err)
	}
	return uc.repo.CreateColumn(*column)
}
