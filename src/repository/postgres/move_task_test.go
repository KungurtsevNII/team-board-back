package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveTaskColumn(t *testing.T) {
	tests := []struct {
		name        string
		taskID      uuid.UUID
		columnID    uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedErr error
	}{
		{
			name:     "успешное перемещение задачи",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name:     "задача не найдена - 0 строк обновлено",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: ErrTaskNotFoundOrDeleted,
		},
		{
			name:     "задача удалена - 0 строк обновлено",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: ErrTaskNotFoundOrDeleted,
		},
		{
			name:     "нарушение внешнего ключа - колонка не существует",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnError(&pgconn.PgError{Code: "23503", Message: "foreign key violation"})
			},
			expectedErr: &pgconn.PgError{Code: "23503"},
		},
		{
			name:     "перемещение в ту же колонку",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name:     "общая ошибка БД",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnError(errors.New("database connection lost"))
			},
			expectedErr: errors.New("database connection lost"),
		},
		{
			name:     "ошибка подключения",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnError(errors.New("connection refused"))
			},
			expectedErr: errors.New("connection refused"),
		},
		{
			name:     "timeout при обновлении",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectedErr: context.DeadlineExceeded,
		},
		{
			name:     "cancelled context",
			taskID:   uuid.New(),
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "tasks" SET`).
					WillReturnError(context.Canceled)
			},
			expectedErr: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)

			repo := &Repository{pool: mock}
			err = repo.MoveTaskColumn(context.Background(), tt.taskID, tt.columnID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				if pgErr, ok := tt.expectedErr.(*pgconn.PgError); ok {
					assert.Contains(t, err.Error(), pgErr.Code)
				} else {
					assert.ErrorContains(t, err, tt.expectedErr.Error())
				}
			} else {
				require.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
