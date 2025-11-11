package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastOrderNumColumn(t *testing.T) {
	tests := []struct {
		name        string
		boardID     uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expectedNum int64
		expectedErr bool
	}{
		{
			name:    "успешное получение последнего order_num",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(5))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedNum: 5,
			expectedErr: false,
		},
		{
			name:    "нет колонок в доске",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnError(pgx.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedNum: 0,
			expectedErr: true,
		},
		{
			name:    "ошибка при начале транзакции",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			expectedNum: 0,
			expectedErr: true,
		},
		{
			name:    "ошибка при коммите",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(3))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			},
			expectedNum: 0,
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
			orderNum, err := repo.GetLastOrderNumColumn(context.Background(), tt.boardID)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNum, orderNum)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
