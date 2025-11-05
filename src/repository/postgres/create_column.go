package postgres

import (
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	// "github.com/jackc/pgx/v5/pgxpool"
)

type ColumnRecord struct {
	ID        string     `db:"id"`
	BoardID   string     `db:"board_id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"primaryKey"`
	UpdatedAt time.Time  `db:"primaryKey"`
}

// func (r Repository) Check(){
// }

func (r Repository) CreateColumn(column domain.Column) error {
	// column domain -> column record
	// record раскаладывать в SQL
	// SQL запрос отправить
	return nil
}
