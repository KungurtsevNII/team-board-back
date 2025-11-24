package postgres

import (
	"context"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r Repository) GetTasks(ctx context.Context, ID string) ([]domain.Task, error) {
	const op = "postgres.GetTasks"
	uid, err := uuid.Parse(ID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	tasks := make([]domain.Task, 0)
	err = pgxscan.Select(ctx, r.pool, tasks,
		`SELECT id, column_id, board_id, number, title 
		FROM tasks WHERE board_id = $1 
		ORDER BY number;`, uid)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return tasks, nil
}
