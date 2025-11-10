package getboard

import (
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
)

type Repo interface {
	GetBoard(ID string) (domain.Board, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(cmd GetBoardCommand) (domain.Board, error) {
	const op = "getboard.Handle"

	board, err := uc.repo.GetBoard(cmd.ID)
	if err != nil {
		return domain.Board{}, fmt.Errorf("%s: %v", op, err)
	}
	return board, nil
}
