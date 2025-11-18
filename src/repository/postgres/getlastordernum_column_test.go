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
		expectedErr error
	}{
		{
			name:    "успешное получение последнего order_num",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(5))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
			},
			expectedNum: 5,
			expectedErr: nil,
		},
		{
			name:    "нет колонок в доске - ErrNoRows",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedNum: 0,
			expectedErr: pgx.ErrNoRows,
		},
		{
			name:    "order_num равен нулю",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(0))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
			},
			expectedNum: 0,
			expectedErr: nil,
		},
		{
			name:    "большое значение order_num",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(999999))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
			},
			expectedNum: 999999,
			expectedErr: nil,
		},
		{
			name:    "отрицательное значение order_num",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow(int64(-1))
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
			},
			expectedNum: -1,
			expectedErr: nil,
		},
		{
			name:    "общая ошибка БД",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnError(errors.New("database error"))
			},
			expectedNum: 0,
			expectedErr: errors.New("database error"),
		},
		{
			name:    "ошибка сканирования - неверный тип данных",
			boardID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"order_num"}).AddRow("not_a_number")
				mock.ExpectQuery(`SELECT "order_num" FROM "columns"`).
					WillReturnRows(rows)
			},
			expectedNum: 0,
			expectedErr: errors.New("'int64' not supported for value kind 'string'"),
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

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.Equal(t, int64(0), orderNum)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNum, orderNum)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
