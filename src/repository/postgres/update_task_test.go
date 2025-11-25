package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

func TestUpdateTask(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name        string
		task        *domain.Task
		mockSetup   func(mock pgxmock.PgxPoolIface, task *domain.Task)
		expectedErr error
	}{
		{
			name: "успешное обновление задачи",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "Updated Task",
				Description: nil,
				Tags:        []string{"updated", "tag1"},
				Checklists:  []domain.Checklist{{Title: "Updated Checklist"}},
				UpdatedAt:   now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление с nil полями",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      5,
				Title:       "Minimal Update",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				UpdatedAt:   now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление с множественными тегами и чеклистами",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      10,
				Title:       "Task with details",
				Description: nil,
				Tags:        []string{"tag1", "tag2", "tag3"},
				Checklists: []domain.Checklist{
					{
						Title: "Checklist 1",
						Items: []domain.ChecklistItem{{Title: "item1", Completed: false}},
					},
				},
				UpdatedAt: now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление несуществующей задачи (0 строк)",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      999,
				Title:       "Nonexistent Update",
				Description: nil,
				Tags:        nil,
				Checklists:  nil,
				UpdatedAt:   now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: nil,
		},
		{
			name: "нарушение внешнего ключа при обновлении",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(), // Invalid FK
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "FK Violation",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				UpdatedAt:   now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnError(pgErr)
			},
			expectedErr: errors.New("foreign key violation"), // Простая строка вместо *PgError
		},
		{
			name: "общая ошибка БД",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "DB Error Update",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				UpdatedAt:   now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				mock.ExpectExec(`UPDATE "tasks"`).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedErr: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.task)

			repo := &Repository{pool: mock}
			err = repo.UpdateTask(context.Background(), tt.task)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				
				if strings.Contains(tt.name, "внешнего ключа") {
					assert.Contains(t, err.Error(), "SQLSTATE 23503")
				}
				
				if !strings.Contains(tt.name, "JSON") {
					assert.NoError(t, mock.ExpectationsWereMet())
				}
			} else {
				require.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
