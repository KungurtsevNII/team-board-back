package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateColumn(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		column      *domain.Column
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedErr error
	}{
		{
			name: "успешное обновление колонки",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "In Progress",
				OrderNum:  2,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "columns"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление колонки с DeletedAt",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Done",
				OrderNum:  3,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: &now,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "columns"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление с пустым именем",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "",
				OrderNum:  0,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "columns"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			expectedErr: nil,
		},
		{
			name: "обновление несуществующей колонки (0 строк)",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Nonexistent",
				OrderNum:  5,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// функция не проверяет rowsAffected, поэтому это считается успехом
				mock.ExpectExec(`UPDATE "columns"`).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			expectedErr: nil,
		},
		{
			name: "ошибка БД при обновлении",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "DB Error",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(`UPDATE "columns"`).
					WillReturnError(errors.New("update error"))
			},
			expectedErr: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)

			repo := &Repository{pool: mock}
			err = repo.UpdateColumn(context.Background(), tt.column)

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
