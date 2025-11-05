package createcolumn

import (
	"github.com/KungurtsevNII/team-board-back/src/domain"
)

type UC struct {
	repo Repository
}

type Repo interface {
	Check()
	Create(column domain.Column) error
}

func (uc *UC) Handle(cmd CreateColumnCommand) error {
	return uc.repo.Create(column)
}
