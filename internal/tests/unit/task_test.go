package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTaskStore struct {
	mock.Mock
}

func (m *mockTaskStore) Create(ctx context.Context, id, userID, title, description string, isCompleted bool) error {
	args := m.Called(ctx, id, userID, title, description, isCompleted)
	return args.Error(0)
}

func (m *mockTaskStore) Get(ctx context.Context, id, userID string) (*models.TaskModel, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskModel), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name            string
		userID          string
		title           string
		description     string
		setupMock       func(store *mockTaskStore)
		wantErr         error
		wantErrContains string
		wantCreate      bool
		wantGet         bool
	}{
		{
			name:        "success",
			userID:      "user-1",
			title:       "Write tests",
			description: "Add unit tests for CreateTask",
			setupMock: func(store *mockTaskStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "user-1", "Write tests", "Add unit tests for CreateTask", false).
					Return(nil)
				store.
					On("Get", mock.Anything, mock.AnythingOfType("string"), "user-1").
					Return(&models.TaskModel{
						ID:          "generated-id",
						UserID:      "user-1",
						Title:       "Write tests",
						Description: "Add unit tests for CreateTask",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil)
			},
			wantCreate: true,
			wantGet:    true,
		},
		{
			name:        "invalid title",
			userID:      "user-1",
			title:       "   ",
			description: "A valid description",
			setupMock:   func(store *mockTaskStore) {},
			wantErr:     models.ErrInvalidTask,
		},
		{
			name:        "invalid description",
			userID:      "user-1",
			title:       "A valid title",
			description: "   ",
			setupMock:   func(store *mockTaskStore) {},
			wantErr:     models.ErrInvalidTask,
		},
		{
			name:        "whitespace values are accepted",
			userID:      "user-1",
			title:       "  Write tests  ",
			description: "  Add unit tests  ",
			setupMock: func(store *mockTaskStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "user-1", "  Write tests  ", "  Add unit tests  ", false).
					Return(nil)
				store.
					On("Get", mock.Anything, mock.AnythingOfType("string"), "user-1").
					Return(&models.TaskModel{ID: "generated-id", UserID: "user-1", Title: "  Write tests  ", Description: "  Add unit tests  "}, nil)
			},
			wantCreate: true,
			wantGet:    true,
		},
		{
			name:        "store create failure",
			userID:      "user-1",
			title:       "Write tests",
			description: "Add unit tests",
			setupMock: func(store *mockTaskStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "user-1", "Write tests", "Add unit tests", false).
					Return(errors.New("create failed"))
			},
			wantErrContains: "create failed",
			wantCreate:      true,
		},
		{
			name:        "store get failure",
			userID:      "user-1",
			title:       "Write tests",
			description: "Add unit tests",
			setupMock: func(store *mockTaskStore) {
				store.
					On("Create", mock.Anything, mock.AnythingOfType("string"), "user-1", "Write tests", "Add unit tests", false).
					Return(nil)
				store.
					On("Get", mock.Anything, mock.AnythingOfType("string"), "user-1").
					Return((*models.TaskModel)(nil), errors.New("get failed"))
			},
			wantErrContains: "get failed",
			wantCreate:      true,
			wantGet:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockTaskStore{}
			if tc.setupMock != nil {
				tc.setupMock(store)
			}

			service := models.NewTaskService(store)
			task, err := service.CreateTask(ctx, tc.userID, tc.title, tc.description)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, task)
			} else if tc.wantErrContains != "" {
				assert.ErrorContains(t, err, tc.wantErrContains)
				assert.Nil(t, task)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, tc.userID, task.UserID)
				assert.Equal(t, tc.title, task.Title)
				assert.Equal(t, tc.description, task.Description)
				assert.NotEmpty(t, task.ID)
			}

			if tc.wantCreate {
				store.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("string"), tc.userID, tc.title, tc.description, false)
			} else {
				store.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}

			if tc.wantGet {
				store.AssertCalled(t, "Get", mock.Anything, mock.AnythingOfType("string"), tc.userID)
			} else {
				store.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
			}

			store.AssertExpectations(t)
		})
	}
}
