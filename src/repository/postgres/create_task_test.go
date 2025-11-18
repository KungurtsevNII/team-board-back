package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTask(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		task        *domain.Task
		mockSetup   func(mock pgxmock.PgxPoolIface, task *domain.Task)
		expectedErr error
	}{
		{
			name: "успешное создание задачи",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "Test Task",
				Description: stringPtr("Test Description"),
				Tags:        []string{"tag1", "tag2"},
				Checklists:  []domain.Checklist{{Title: "Checklist 1"}},
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание задачи с nil полями",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      5,
				Title:       "Minimal Task",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание задачи с большими тегами",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      10,
				Title:       "Task with many tags",
				Description: stringPtr("Description"),
				Tags:        []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
				Checklists: []domain.Checklist{
					{
						Title: "Checklist 1",
						Items: []domain.ChecklistItem{
							{Title: "item1", Completed: false},
						},
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "создание задачи с удаленной датой",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      3,
				Title:       "Deleted Task",
				Description: stringPtr("Deleted"),
				Tags:        []string{"deleted"},
				Checklists:  nil,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   timePtr(now),
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			expectedErr: nil,
		},
		{
			name: "ошибка БД при вставке",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "Test Task",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnError(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
		{
			name: "нарушение уникального ограничения",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "Duplicate Task",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			expectedErr: &pgconn.PgError{Code: "23505"},
		},
		{
			name: "нарушение внешнего ключа",
			task: &domain.Task{
				ID:          uuid.New(),
				BoardID:     uuid.New(),
				ColumnID:    uuid.New(),
				Number:      1,
				Title:       "Invalid FK",
				Description: nil,
				Tags:        []string{},
				Checklists:  nil,
				CreatedAt:   now,
				UpdatedAt:   now,
				DeletedAt:   nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, task *domain.Task) {
				checklistsJSON, _ := json.Marshal(task.Checklists)
				mock.ExpectExec(`INSERT INTO tasks`).
					WithArgs(
						task.ID,
						task.BoardID,
						task.ColumnID,
						task.Number,
						task.Title,
						task.Description,
						task.Tags,
						checklistsJSON,
						task.CreatedAt,
						task.UpdatedAt,
						task.DeletedAt,
					).
					WillReturnError(&pgconn.PgError{Code: "23503"})
			},
			expectedErr: &pgconn.PgError{Code: "23503"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.task)

			repo := &Repository{pool: mock}
			err = repo.CreateTask(context.Background(), tt.task)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Вспомогательные функции для создания указателей
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
