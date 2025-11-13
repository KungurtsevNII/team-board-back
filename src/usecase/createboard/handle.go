package createboard

import (
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
)

type Repo interface {
	CheckBoard(shortName string) bool
	CreateBoard(board domain.Board) (string, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(cmd CreateBoardCommand) (string, error) {

	const op = "createboard.Handle"
	if uc.repo.CheckBoard(cmd.ShortName) {
		return "", BoardIsExistsErr
	}

	board, err := domain.NewBoard(cmd.Name, cmd.ShortName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := uc.repo.CreateBoard(*board)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
