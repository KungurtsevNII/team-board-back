package postgres

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

    row := tx.QueryRow(ctx, sql, params...)
	err = row.Scan(&orderNum)
	if err != nil{
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    if err := tx.Commit(ctx); err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return orderNum, nil
}

