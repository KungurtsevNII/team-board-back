package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewColumn(t *testing.T) {
	boardID := uuid.New()

	testCases := []struct {
		name        string
		boardID     uuid.UUID
		columnName  string
		orderNum    int64
		expectError bool
		errorType   error
	}{
		{
			name:        "Success: creates a new column with valid data",
			boardID:     boardID,
			columnName:  "To Do",
			orderNum:    1,
			expectError: false,
		},
		{
			name:        "Failure: returns error for empty name",
			boardID:     boardID,
			columnName:  "",
			orderNum:    1,
			expectError: true,
			errorType:   ErrEmptyColumnName,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			column, err := NewColumn(tc.boardID, tc.columnName, tc.orderNum)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.errorType)
				assert.Nil(t, column)
			} else {
				require.NoError(t, err)
				require.NotNil(t, column)
				assert.NotEqual(t, uuid.Nil, column.ID, "ID should be initialized")
				assert.Equal(t, tc.boardID, column.BoardID)
				assert.Equal(t, tc.columnName, column.Name)
				assert.Equal(t, tc.orderNum, column.OrderNum)
				assert.WithinDuration(t, time.Now().UTC(), column.CreatedAt, time.Second, "CreatedAt should be recent")
				assert.WithinDuration(t, time.Now().UTC(), column.UpdatedAt, time.Second, "UpdatedAt should be recent")
				assert.Nil(t, column.DeletedAt, "DeletedAt should be nil on creation")
			}
		})
	}
}

func TestColumn_Delete(t *testing.T) {
	column, err := NewColumn(uuid.New(), "In Progress", 2)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		columnToDelete *Column
	}{
		{
			name:           "marks a column as deleted",
			columnToDelete: column,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.columnToDelete)
			require.Nil(t, tc.columnToDelete.DeletedAt)

			tc.columnToDelete.Delete()

			require.NotNil(t, tc.columnToDelete.DeletedAt, "DeletedAt should be set after delete")
			assert.WithinDuration(t, time.Now().UTC(), *tc.columnToDelete.DeletedAt, time.Second, "DeletedAt should be a recent timestamp")
		})
	}
}
