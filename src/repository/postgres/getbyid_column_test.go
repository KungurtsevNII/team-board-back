package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetColumnByID(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		columnID    uuid.UUID
		mockSetup   func(mock pgxmock.PgxPoolIface, columnID uuid.UUID)
		expected    *domain.Column
		expectedErr error
	}{
		{
			name:     "успешное получение колонки",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				rows := pgxmock.NewRows([]string{
					"id",
					"board_id",
					"name",
					"order_num",
					"created_at",
					"deleted_at",
					"updated_at",
				}).AddRow(
					columnID,
					uuid.MustParse("11111111-1111-1111-1111-111111111111"), // board_id
					"Todo",    // name
					int64(10), // order_num
					now,       // created_at
					nil,       // deleted_at
					now,       // updated_at
				)

				mock.ExpectQuery(`SELECT .* FROM "columns"`).
					WillReturnRows(rows)
			},
			expected: &domain.Column{
				// ID и BoardID дополним ниже
				Name:     "Todo",
				OrderNum: 10,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: zeroTimePtr(),
			},
			expectedErr: nil,
		},
		{
			name:     "колонка не найдена - ErrNoRows",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				mock.ExpectQuery(`SELECT .* FROM "columns"`).
					WillReturnError(pgx.ErrNoRows)
			},
			expected:    nil,
			expectedErr: pgx.ErrNoRows,
		},
		{
			name:     "ошибка БД при выборке",
			columnID: uuid.New(),
			mockSetup: func(mock pgxmock.PgxPoolIface, columnID uuid.UUID) {
				mock.ExpectQuery(`SELECT .* FROM "columns"`).
					WillReturnError(errors.New("database error"))
			},
			expected:    nil,
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockSetup(mock, tt.columnID)

			repo := &Repository{pool: mock}
			col, err := repo.GetColumnByID(context.Background(), tt.columnID)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
				assert.Nil(t, col)
			} else {
				require.NoError(t, err)
				require.NotNil(t, col)

				// Заполним ожидаемые ID, если они не заданы
				if tt.expected.ID == uuid.Nil {
					tt.expected.ID = tt.columnID
				}
				if tt.expected.BoardID == uuid.Nil {
					tt.expected.BoardID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
				}

				assert.Equal(t, tt.expected.ID, col.ID)
				assert.Equal(t, tt.expected.BoardID, col.BoardID)
				assert.Equal(t, tt.expected.Name, col.Name)
				assert.Equal(t, tt.expected.OrderNum, col.OrderNum)
				assert.WithinDuration(t, tt.expected.CreatedAt, col.CreatedAt, time.Second)
				assert.WithinDuration(t, tt.expected.UpdatedAt, col.UpdatedAt, time.Second)
				assert.Equal(t, tt.expected.DeletedAt, col.DeletedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
