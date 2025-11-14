package postgres

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (r Repository) GetLastNumberTask(ctx context.Context, boardID uuid.UUID) (int64, error) {
    const op = "postgres.GetLastNumberTask"

    ds := goqu.From("tasks").
		Where(goqu.C("board_id").Eq(boardID)).
		Select(goqu.C("number")).
		Order(goqu.I("number").Desc()).Limit(1)
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var number int64
    row := r.pool.QueryRow(ctx, sql, params...)
	err = row.Scan(&number)
	if err != nil{
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return number, nil
}

