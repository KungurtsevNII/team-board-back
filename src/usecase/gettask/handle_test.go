package gettask

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/gettask/mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	taskID := uuid.New()
	ctx := context.Background()

	// Вспомогательная функция для создания указателя на строку
	strPtr := func(s string) *string {
		return &s
	}

	expectedTask := &domain.Task{
		ID:          taskID,
		ColumnID:    uuid.New(),
		BoardID:     uuid.New(),
		Number:      1,
		Title:       "Test Task",
		Description: strPtr("A description"),
		Tags:        []string{"test"},
		Checklists:  []domain.Checklist{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	testCases := []struct {
		name        string
		query       GetTaskQuery
		setupMock   func(*mocks.Repo)
		expected    *domain.Task
		expectError error
	}{
		{
			name:  "Success: task found",
			query: GetTaskQuery{TaskID: taskID},
			setupMock: func(repo *mocks.Repo) {
				repo.On("GetTaskByID", mock.Anything, taskID).
					Return(expectedTask, nil).
					Once()
			},
			expected:    expectedTask,
			expectError: nil,
		},
		{
			name:  "Failure: task not found",
			query: GetTaskQuery{TaskID: taskID},
			setupMock: func(repo *mocks.Repo) {
				repo.On("GetTaskByID", mock.Anything, taskID).
					Return(nil, pgx.ErrNoRows).
					Once()
			},
			expected:    nil,
			expectError: ErrTaskNotFound,
		},
		{
			name:  "Failure: unknown repository error",
			query: GetTaskQuery{TaskID: taskID},
			setupMock: func(repo *mocks.Repo) {
				repo.On("GetTaskByID", mock.Anything, taskID).
					Return(nil, errors.New("unexpected db error")). 
					Once()
			},
			expected:    nil,
			expectError: ErrGetTaskUnknown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepo(t)
			tc.setupMock(repo)

			uc := NewUC(repo)

			task, err := uc.Handle(ctx, tc.query)

			if tc.expectError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectError, "Wrong error type")
				assert.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
				assert.Equal(t, tc.expected, task)
			}

			repo.AssertExpectations(t)
		})
	}
}
