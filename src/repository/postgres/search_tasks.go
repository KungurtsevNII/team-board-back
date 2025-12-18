package postgres

import (
	"github.com/KungurtsevNII/team-board-back/src/domain"
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lib/pq"
)

func (r Repository) SearchTasks(
    ctx context.Context,
    tags []string,
    query string,
    limit, offset uint,
) ([]domain.Task, error) {
    const op = "postgres.SearchTasks"

    tasksRec := make([]TaskSearchRecord, 0)
    
    ds := goqu.From("tasks")
    
    if len(tags) > 0 {
        ds = ds.Where(goqu.L("tags @> ?", pq.Array(tags)))
    }
    
    if query != "" {
        ds = ds.Where(goqu.Ex{"title": goqu.Op{"ilike": "%" + query + "%"}})
    }
    
    ds = ds.Select(&TaskSearchRecord{}).
        Join(goqu.T("boards"), goqu.On(goqu.T("tasks").Col("board_id").Eq(goqu.T("boards").Col("id")))).
        Join(goqu.T("columns"), goqu.On(goqu.T("tasks").Col("column_id").Eq(goqu.T("columns").Col("id")))).
        Where(goqu.T("tasks").Col("deleted_at").IsNull(), 
            goqu.T("boards").Col("deleted_at").IsNull()).
		Order(goqu.T("tasks").Col("created_at").Desc()).
        Limit(limit).
        Offset(offset)
    
    sql, params, err := ds.ToSQL()
    if err != nil {
        return nil, errors.Wrap(err, op)
    }
    
    err = pgxscan.Select(ctx, r.pool, &tasksRec, sql, params...)
    if err != nil {
        return nil, errors.Wrap(err, op)
    }
    
    dmn, err := TaskSearchRecords(tasksRec).toDomain()
    if err != nil {
        return nil, errors.Wrap(err, op)
    }
    
    return dmn, nil
}
