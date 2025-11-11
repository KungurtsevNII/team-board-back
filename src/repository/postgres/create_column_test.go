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

func TestCreateColumn(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		column      *domain.Column
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedErr bool
	}{
		{
			name: "успешное создание колонки",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "To Do",
				OrderNum:  1,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
			expectedErr: false,
		},
		{
			name: "ошибка при начале транзакции",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "In Progress",
				OrderNum:  3,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New("begin tx error"))
			},
			expectedErr: true,
		},
		{
			name: "ошибка при INSERT",
			column: &domain.Column{
				ID:        uuid.New(),
				BoardID:   uuid.New(),
				Name:      "Review",
				OrderNum:  4,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "columns"`).
					WillReturnError(errors.New("constraint violation"))
				mock.ExpectRollback()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock)

			repo := &Repository{pool: mock}
			err = repo.CreateColumn(context.Background(), tt.column)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
