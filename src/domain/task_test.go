package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string {
	return &s
}

func TestNewTask(t *testing.T) {
	t.Run("Success: creates a new task with valid data", func(t *testing.T) {
		columnID := uuid.New()
		boardID := uuid.New()
		number := int64(101)
		title := "Implement API"
		description := strPtr("Use Golang")
		tags := []string{"backend", "api"}
		checklists := []Checklist{}

		task, err := NewTask(columnID, boardID, number, title, description, tags, checklists)

		require.NoError(t, err)
		require.NotNil(t, task)

		assert.NotEqual(t, uuid.Nil, task.ID)
		assert.Equal(t, columnID, task.ColumnID)
		assert.Equal(t, boardID, task.BoardID)
		assert.Equal(t, number, task.Number)
		assert.Equal(t, title, task.Title)
		assert.Equal(t, description, task.Description)
		assert.Equal(t, tags, task.Tags)
		assert.Equal(t, checklists, task.Checklists)
		assert.WithinDuration(t, time.Now().UTC(), task.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now().UTC(), task.UpdatedAt, time.Second)
		assert.Nil(t, task.DeletedAt)
	})
}

func TestTask_Update(t *testing.T) {
	t.Run("updates task fields and timestamp", func(t *testing.T) {
		task, err := NewTask(uuid.New(), uuid.New(), 1, "Old Title", nil, nil, nil)
		require.NoError(t, err)
		originalUpdatedAt := task.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		newColumnID := uuid.New()
		newBoardID := uuid.New()
		newNumber := int64(2)
		newTitle := "New Title"
		newDescription := strPtr("New Desc")
		newTags := []string{"new"}
		newChecklists := []Checklist{{Title: "new list"}}

		task.Update(newColumnID, newBoardID, newNumber, newTitle, newDescription, newTags, newChecklists)

		assert.Equal(t, newColumnID, task.ColumnID)
		assert.Equal(t, newBoardID, task.BoardID)
		assert.Equal(t, newNumber, task.Number)
		assert.Equal(t, newTitle, task.Title)
		assert.Equal(t, newDescription, task.Description)
		assert.Equal(t, newTags, task.Tags)
		assert.Equal(t, newChecklists, task.Checklists)
		assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
	})
}

func TestTask_Delete(t *testing.T) {
	t.Run("marks a task as deleted", func(t *testing.T) {
		task, err := NewTask(uuid.New(), uuid.New(), 1, "A task", nil, nil, nil)
		require.NoError(t, err)
		require.Nil(t, task.DeletedAt)

		task.Delete()

		require.NotNil(t, task.DeletedAt)
		assert.WithinDuration(t, time.Now().UTC(), *task.DeletedAt, time.Second)
	})
}

func TestTask_MoveToColumn(t *testing.T) {
	t.Run("Success: moves task to a different column", func(t *testing.T) {
		originalColumnID := uuid.New()
		task, err := NewTask(originalColumnID, uuid.New(), 1, "Task to move", nil, nil, nil)
		require.NoError(t, err)
		originalUpdatedAt := task.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		newColumnID := uuid.New()

		err = task.MoveToColumn(newColumnID)

		require.NoError(t, err)
		assert.Equal(t, newColumnID, task.ColumnID)
		assert.True(t, task.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("Failure: returns error when moving to the same column", func(t *testing.T) {
		originalColumnID := uuid.New()
		task, err := NewTask(originalColumnID, uuid.New(), 1, "Task not to move", nil, nil, nil)
		require.NoError(t, err)
		originalUpdatedAt := task.UpdatedAt

		err = task.MoveToColumn(originalColumnID) 

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrAlreadyInColumn)
		assert.Equal(t, originalColumnID, task.ColumnID)
		assert.Equal(t, originalUpdatedAt, task.UpdatedAt, "UpdatedAt should not change")
	})
}
