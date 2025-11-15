package postgres

import (
	// "github.com/jackc/pgx/v5/pgxpool"
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (r Repository) CheckColumnInBoard(ctx context.Context, boardID uuid.UUID, columnID uuid.UUID) (bool, error) {
    const op = "postgres.CheckColumnInBoard"
    
    ds := goqu.From("columns").
    Select(goqu.COUNT("*")).
    Where(
        goqu.C("board_id").Eq(boardID), 
        goqu.C("id").Eq(columnID),
    )
    
    sql, params, err := ds.ToSQL()
    if err != nil {
        return false, fmt.Errorf("%s: %w", op, err)
    }
    
    var count int64
    err = r.pool.QueryRow(ctx, sql, params...).Scan(&count)
    if err != nil {
        return false, fmt.Errorf("%s: %w", op, err)
    }
    
    return count > 0, nil
}
