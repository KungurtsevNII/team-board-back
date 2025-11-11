package postgres

import (
	"context"
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
)

func (r Repository) CreateColumn(ctx context.Context, column *domain.Column) error {
    const op = "postgres.CreateColumn"

	record := ColumnRecord{
		ID: column.ID,
		BoardID: column.BoardID,
		Name: column.Name,
		OrderNum: column.OrderNum,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
		DeletedAt: column.DeletedAt,
	}

    ds := goqu.Insert("columns").
        Rows(record).
        Returning(goqu.C("id"))

    sql, params, err := ds.ToSQL()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }


    tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    defer tx.Rollback(ctx)

    _, err = tx.Exec(ctx, sql, params...)
	if err != nil{
        return fmt.Errorf("%s: %w", op, err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    return nil
}

