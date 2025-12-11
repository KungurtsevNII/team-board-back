package postgres

import (
	"context"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
)

func (r Repository) UpdateBoard(ctx context.Context, board *domain.Board) error {
    const op = "postgres.UpdateBoard"
    board.UpdatedAt = time.Now().UTC()

    // Тут сделал транзакцию, чтобы избежать фигни всякой
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return errors.Wrap(err, op)
    }
    defer tx.Rollback(ctx)

    ds := goqu.Update("boards").Where(
        goqu.C("id").Eq(board.ID),
        goqu.C("deleted_at").IsNull(),
    ).Set(
        BoardRecord{
            Name:      board.Name,
            ShortName: board.ShortName,
            UpdatedAt: board.UpdatedAt,
            DeletedAt: board.DeletedAt,
        },
    )

    sql, params, err := ds.ToSQL()
    if err != nil {
        return errors.Wrap(err, op)
    }

    _, err = tx.Exec(ctx, sql, params...)
    if err != nil {
        return errors.Wrap(err, op)
    }

    if board.DeletedAt != nil {
        dsTasks := goqu.Update("tasks").
            Where(
                goqu.C("board_id").Eq(board.ID),
                goqu.C("deleted_at").IsNull(),
            ).
            Set(goqu.Record{
                "deleted_at": board.DeletedAt,
                "updated_at": board.UpdatedAt,
            })
        
        sqlTasks, paramsTasks, err := dsTasks.ToSQL()
        if err != nil {
            return errors.Wrap(err, op)
        }
        
        if _, err := tx.Exec(ctx, sqlTasks, paramsTasks...); err != nil {
            return errors.Wrap(err, op)
        }

        dsColumns := goqu.Update("columns").
            Where(
                goqu.C("board_id").Eq(board.ID),
                goqu.C("deleted_at").IsNull(),
            ).
            Set(goqu.Record{
                "deleted_at": board.DeletedAt,
                "updated_at": board.UpdatedAt,
            })

        sqlCols, paramsCols, err := dsColumns.ToSQL()
        if err != nil {
            return errors.Wrap(err, op)
        }

        if _, err := tx.Exec(ctx, sqlCols, paramsCols...); err != nil {
            return errors.Wrap(err, op)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        return errors.Wrap(err, op)
    }

    return nil
}

