package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewColumn(t *testing.T) {
	t.Run("Success: creates a new column with valid data", func(t *testing.T) {
		boardID := uuid.New()
		name := "To Do"
		orderNum := int64(1)

		column, err := NewColumn(boardID, name, orderNum)

		require.NoError(t, err)
		require.NotNil(t, column)

		assert.NotEqual(t, uuid.Nil, column.ID, "ID should be initialized")
		assert.Equal(t, boardID, column.BoardID)
		assert.Equal(t, name, column.Name)
		assert.Equal(t, orderNum, column.OrderNum)
		assert.WithinDuration(t, time.Now().UTC(), column.CreatedAt, time.Second, "CreatedAt should be recent")
		assert.WithinDuration(t, time.Now().UTC(), column.UpdatedAt, time.Second, "UpdatedAt should be recent")
		assert.Nil(t, column.DeletedAt, "DeletedAt should be nil on creation")
	})

	t.Run("Failure: returns error for empty name", func(t *testing.T) {
		boardID := uuid.New()
		emptyName := ""
		orderNum := int64(1)

		column, err := NewColumn(boardID, emptyName, orderNum)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrEmptyColumnName)
		assert.Nil(t, column)
	})
}

func TestColumn_Delete(t *testing.T) {
	t.Run("marks a column as deleted", func(t *testing.T) {
		column, err := NewColumn(uuid.New(), "In Progress", 2)
		require.NoError(t, err)
		require.NotNil(t, column)
		require.Nil(t, column.DeletedAt)

		column.Delete()

		require.NotNil(t, column.DeletedAt, "DeletedAt should be set after delete")
		assert.WithinDuration(t, time.Now().UTC(), *column.DeletedAt, time.Second, "DeletedAt should be a recent timestamp")
	})
}
