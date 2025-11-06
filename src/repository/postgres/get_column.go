package postgres

import (
	"github.com/KungurtsevNII/team-board-back/src/domain"
	// "github.com/jackc/pgx/v5/pgxpool"
)


func (r Repository) GetColumn(ID string) (domain.Column, error) {
	// column domain -> column record
	// record раскаладывать в SQL
	// SQL запрос отправить
	return domain.Column{}, nil
}
