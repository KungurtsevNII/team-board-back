package postgres

import (
	"context"
	"fmt"

	"github.com/KungurtsevNII/team-board-back/src/domain"
)

func (r Repository) CreateTask(ctx context.Context, task *domain.Task) error {
    const op = "postgres.CreateTask"
    
	//Гоку не хочет адекватно жрать теги []string и чеклисты json.RawMessage, поэтому пришлось вручную
    sql := `INSERT INTO tasks (
				id, board_id, column_id, number, title, description, tags, 
				checklists, created_at, updated_at, deleted_at
			)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
    
    _, err := r.pool.Exec(ctx, sql,
        task.ID,
        task.BoardID,
        task.ColumnID,
        task.Number,
        task.Title,
        task.Description,
        task.Tags,
        task.Checklists,
        task.CreatedAt,
        task.UpdatedAt,
        task.DeletedAt,
    )
    
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }
    
    return nil
}


