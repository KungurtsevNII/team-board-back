package postgres

import "github.com/KungurtsevNII/team-board-back/src/domain"

func (r Repository) CreateBoard(board domain.Board) (string, error) {
	// board domain -> board record
	// record раскаладывать в SQL
	// SQL запрос отправить
	return "1", nil
}
