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

func (m *mockTaskStore) Get(ctx context.Context, id, userID string) (*models.Task, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *mockTaskStore) Update(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *mockTaskStore) List(ctx context.Context, userID string) ([]models.Task, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *mockTaskStore) Delete(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
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
					Return(&models.Task{
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
					Return(&models.Task{ID: "generated-id", UserID: "user-1", Title: "  Write tests  ", Description: "  Add unit tests  "}, nil)
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
					Return((*models.Task)(nil), errors.New("get failed"))
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

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()

	newString := func(value string) *string { return &value }
	newBool := func(value bool) *bool { return &value }

	cases := []struct {
		name            string
		title           *string
		description     *string
		isCompleted     *bool
		setupMock       func(store *mockTaskStore)
		wantErr         error
		wantErrContains string
		wantUpdate      bool
		wantUpdatedTask func(t *testing.T, task models.TaskModel)
	}{
		{
			name:        "success updates all fields",
			title:       newString("Updated title"),
			description: newString("Updated description"),
			isCompleted: newBool(true),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task != nil && task.ID == "task-1" && task.UserID == "user-1" && task.Title == "Updated title" && task.Description == "Updated description" && task.IsCompleted
				})).Return(nil)
			},
			wantUpdate: true,
			wantUpdatedTask: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.Equal(t, true, task.IsCompleted)
			},
		},
		{
			name:        "update only title",
			title:       newString("Updated title"),
			description: nil,
			isCompleted: nil,
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task != nil && task.Title == "Updated title" && task.Description == "Original description" && !task.IsCompleted
				})).Return(nil)
			},
			wantUpdate: true,
			wantUpdatedTask: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.Equal(t, false, task.IsCompleted)
			},
		},
		{
			name:        "update only description",
			title:       nil,
			description: newString("Updated description"),
			isCompleted: nil,
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task != nil && task.Title == "Original title" && task.Description == "Updated description" && !task.IsCompleted
				})).Return(nil)
			},
			wantUpdate: true,
			wantUpdatedTask: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.Equal(t, false, task.IsCompleted)
			},
		},
		{
			name:        "update only completion status",
			title:       nil,
			description: nil,
			isCompleted: newBool(true),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task != nil && task.Title == "Original title" && task.Description == "Original description" && task.IsCompleted
				})).Return(nil)
			},
			wantUpdate: true,
			wantUpdatedTask: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.Equal(t, true, task.IsCompleted)
			},
		},
		{
			name:        "no-op update keeps the task unchanged",
			title:       nil,
			description: nil,
			isCompleted: nil,
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.MatchedBy(func(task *models.Task) bool {
					return task != nil && task.Title == "Original title" && task.Description == "Original description" && !task.IsCompleted
				})).Return(nil)
			},
			wantUpdate: true,
			wantUpdatedTask: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.Equal(t, false, task.IsCompleted)
			},
		},
		{
			name:        "invalid title is rejected",
			title:       newString("   "),
			description: newString("Updated description"),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
			},
			wantErr: models.ErrInvalidTask,
		},
		{
			name:        "invalid description is rejected",
			title:       newString("Updated title"),
			description: newString("   "),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
			},
			wantErr: models.ErrInvalidTask,
		},
		{
			name:        "store get failure",
			title:       newString("Updated title"),
			description: newString("Updated description"),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return((*models.Task)(nil), errors.New("get failed"))
			},
			wantErrContains: "get failed",
		},
		{
			name:        "store update failure",
			title:       newString("Updated title"),
			description: newString("Updated description"),
			isCompleted: newBool(true),
			setupMock: func(store *mockTaskStore) {
				store.On("Get", mock.Anything, "task-1", "user-1").Return(&models.Task{ID: "task-1", UserID: "user-1", Title: "Original title", Description: "Original description"}, nil)
				store.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed"))
			},
			wantErrContains: "update failed",
			wantUpdate:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockTaskStore{}
			if tc.setupMock != nil {
				tc.setupMock(store)
			}

			service := models.NewTaskService(store)
			task, err := service.UpdateTask(ctx, "task-1", "user-1", tc.title, tc.description, tc.isCompleted)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else if tc.wantErrContains != "" {
				assert.ErrorContains(t, err, tc.wantErrContains)
			} else {
				assert.NoError(t, err)
			}

			store.AssertCalled(t, "Get", mock.Anything, "task-1", "user-1")
			if tc.wantUpdate {
				store.AssertCalled(t, "Update", mock.Anything, mock.Anything)
			} else {
				store.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			}

			if tc.wantUpdatedTask != nil {
				tc.wantUpdatedTask(t, *task)
			}

			store.AssertExpectations(t)
		})
	}
}

func TestListTasks(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name            string
		userID          string
		setupMock       func(store *mockTaskStore)
		wantErr         error
		wantErrContains string
		wantTasks       []models.TaskModel
	}{
		{
			name:   "success returns tasks",
			userID: "user-1",
			setupMock: func(store *mockTaskStore) {
				store.On("List", mock.Anything, "user-1").Return([]models.Task{
					{ID: "task-1", UserID: "user-1", Title: "Title 1", Description: "Desc 1", IsCompleted: false},
					{ID: "task-2", UserID: "user-1", Title: "Title 2", Description: "Desc 2", IsCompleted: true},
				}, nil)
			},
			wantTasks: []models.TaskModel{
				{ID: "task-1", UserID: "user-1", Title: "Title 1", Description: "Desc 1", IsCompleted: false},
				{ID: "task-2", UserID: "user-1", Title: "Title 2", Description: "Desc 2", IsCompleted: true},
			},
		},
		{
			name:   "success returns empty list",
			userID: "user-1",
			setupMock: func(store *mockTaskStore) {
				store.On("List", mock.Anything, "user-1").Return([]models.Task{}, nil)
			},
			wantTasks: []models.TaskModel{},
		},
		{
			name:   "store list failure",
			userID: "user-1",
			setupMock: func(store *mockTaskStore) {
				store.On("List", mock.Anything, "user-1").Return([]models.Task(nil), errors.New("list failed"))
			},
			wantErrContains: "list failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &mockTaskStore{}
			if tc.setupMock != nil {
				tc.setupMock(store)
			}

			service := models.NewTaskService(store)
			tasks, err := service.ListTasks(ctx, tc.userID)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, tasks)
			} else if tc.wantErrContains != "" {
				assert.ErrorContains(t, err, tc.wantErrContains)
				assert.Nil(t, tasks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantTasks, tasks)
			}

			store.AssertCalled(t, "List", mock.Anything, tc.userID)
			store.AssertExpectations(t)
		})
	}
}
