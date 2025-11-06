package getcolumn

import (
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
	GetColumn(ID string) (domain.Column, error)
}

func (uc *UC) Handle(cmd GetColumnCommand) (domain.Column, error) {
	//тут прям обращение к базе данных
	const op = "getcolumn.Handle"

	return uc.repo.GetColumn(cmd.ID)
}
