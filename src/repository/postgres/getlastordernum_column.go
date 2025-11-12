package postgres

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (r Repository) GetLastOrderNumColumn(ctx context.Context, boardID uuid.UUID) (orderNum int64, err error) {
    const op = "postgres.GetLastOrderNumColumn"

    ds := goqu.From("columns").
		Where(goqu.C("board_id").Eq(boardID)).
		Select(goqu.C("order_num")).
		Order(goqu.I("order_num").Desc()).Limit(1)
		
	sql, params, err := ds.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

    row := r.pool.QueryRow(ctx, sql, params...)
	err = row.Scan(&orderNum)
	if err != nil{
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return orderNum, nil
}

