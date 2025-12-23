package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	columnID := uuid.New()
	boardID := uuid.New()

	testCases := []struct {
		name        string
		columnID    uuid.UUID
		boardID     uuid.UUID
		number      int64
		title       string
		description *string
		tags        []string
		checklists  []Checklist
	}{
		{
			name:        "Success: creates a new task with valid data",
			columnID:    columnID,
			boardID:     boardID,
			number:      101,
			title:       "Implement API",
			description: strPtr("Use Golang"),
			tags:        []string{"backend", "api"},
			checklists:  []Checklist{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.columnID, tc.boardID, tc.number, tc.title, tc.description, tc.tags, tc.checklists)

			require.NoError(t, err)
			require.NotNil(t, task)

			assert.NotEqual(t, uuid.Nil, task.ID)
			assert.Equal(t, tc.columnID, task.ColumnID)
			assert.Equal(t, tc.boardID, task.BoardID)
			assert.Equal(t, tc.number, task.Number)
			assert.Equal(t, tc.title, task.Title)
			assert.Equal(t, tc.description, task.Description)
			assert.Equal(t, tc.tags, task.Tags)
			assert.Equal(t, tc.checklists, task.Checklists)
			assert.WithinDuration(t, time.Now().UTC(), task.CreatedAt, time.Second)
			assert.WithinDuration(t, time.Now().UTC(), task.UpdatedAt, time.Second)
			assert.Nil(t, task.DeletedAt)
		})
	}
}

func TestTask_Update(t *testing.T) {
	task, err := NewTask(uuid.New(), uuid.New(), 1, "Old Title", nil, nil, nil)
	require.NoError(t, err)
	originalUpdatedAt := task.UpdatedAt

	testCases := []struct {
		name          string
		taskToUpdate  *Task
		newColumnID   uuid.UUID
		newBoardID    uuid.UUID
		newNumber     int64
		newTitle      string
		newDescription *string
		newTags       []string
		newChecklists []Checklist
	}{
		{
			name:          "updates task fields and timestamp",
			taskToUpdate:  task,
			newColumnID:   uuid.New(),
			newBoardID:    uuid.New(),
			newNumber:     2,
			newTitle:      "New Title",
			newDescription: strPtr("New Desc"),
			newTags:       []string{"new"},
			newChecklists: []Checklist{{Title: "new list"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.taskToUpdate.Update(tc.newColumnID, tc.newBoardID, tc.newNumber, tc.newTitle, tc.newDescription, tc.newTags, tc.newChecklists)

			assert.Equal(t, tc.newColumnID, tc.taskToUpdate.ColumnID)
			assert.Equal(t, tc.newBoardID, tc.taskToUpdate.BoardID)
			assert.Equal(t, tc.newNumber, tc.taskToUpdate.Number)
			assert.Equal(t, tc.newTitle, tc.taskToUpdate.Title)
			assert.Equal(t, tc.newDescription, tc.taskToUpdate.Description)
			assert.Equal(t, tc.newTags, tc.taskToUpdate.Tags)
			assert.Equal(t, tc.newChecklists, tc.taskToUpdate.Checklists)
			assert.True(t, tc.taskToUpdate.UpdatedAt.After(originalUpdatedAt))
		})
	}
}

func TestTask_Delete(t *testing.T) {
	task, err := NewTask(uuid.New(), uuid.New(), 1, "A task", nil, nil, nil)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		taskToDelete *Task
	}{
		{
			name:         "marks a task as deleted",
			taskToDelete: task,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Nil(t, tc.taskToDelete.DeletedAt)

			tc.taskToDelete.Delete()

			require.NotNil(t, tc.taskToDelete.DeletedAt)
			assert.WithinDuration(t, time.Now().UTC(), *tc.taskToDelete.DeletedAt, time.Second)
		})
	}
}

func TestTask_MoveToColumn(t *testing.T) {
	originalColumnID := uuid.New()
	taskToMove, err := NewTask(originalColumnID, uuid.New(), 1, "Task to move", nil, nil, nil)
	require.NoError(t, err)

	taskNotToMove, err := NewTask(originalColumnID, uuid.New(), 1, "Task not to move", nil, nil, nil)
	require.NoError(t, err)

	newColumnID := uuid.New()

	testCases := []struct {
		name              string
		task              *Task
		targetColumnID    uuid.UUID
		expectError       bool
		expectedUpdatedAt time.Time
	}{
		{
			name:           "Success: moves task to a different column",
			task:           taskToMove,
			targetColumnID: newColumnID,
			expectError:    false,
		},
		{
			name:              "Failure: returns error when moving to the same column",
			task:              taskNotToMove,
			targetColumnID:    originalColumnID,
			expectError:       true,
			expectedUpdatedAt: taskNotToMove.UpdatedAt, 
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalUpdatedAt := tc.task.UpdatedAt

			err := tc.task.MoveToColumn(tc.targetColumnID)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrAlreadyInColumn)
				assert.Equal(t, tc.targetColumnID, tc.task.ColumnID, "ColumnID should not change on error")
				assert.Equal(t, tc.expectedUpdatedAt, tc.task.UpdatedAt, "UpdatedAt should not change on error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.targetColumnID, tc.task.ColumnID)
				assert.True(t, tc.task.UpdatedAt.After(originalUpdatedAt))
			}
		})
	}
}
