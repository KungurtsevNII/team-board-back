package searchtasks

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/pkg/errors"
)

type Repo interface {
	SearchTasks(
		ctx context.Context, 
		tags []string, 
		query string,
		limit, offset uint) ([]domain.Task, error)
}

type UC struct {
	repo Repo
}

func NewUC(repo Repo) *UC {
	return &UC{
		repo: repo,
	}
}

func (uc *UC) Handle(ctx context.Context, q Query) ([]domain.Task, error) {
	tasks, err := uc.repo.SearchTasks(ctx, q.Tags, q.Query, q.Limit, q.Offset)
	if err != nil {
		return nil, errors.Wrap(ErrSearchTasks, err.Error())
	}

	return tasks, nil
}
