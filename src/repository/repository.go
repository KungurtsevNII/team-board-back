package repository

import "github.com/KungurtsevNII/team-board-back/src/domain"

type RepositoryInf interface {
	CreateColumn(column domain.Column) error
}
