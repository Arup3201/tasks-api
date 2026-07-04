package integration

import (
	"context"
	"testing"

	"github.com/Arup3201/gotask/internal/models"
	"github.com/Arup3201/gotask/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UpdateTaskSuite struct {
	suite.Suite
	db      *gorm.DB
	service *models.TaskService
}

func TestUpdateTaskSuite(t *testing.T) {
	suite.Run(t, new(UpdateTaskSuite))
}

func (s *UpdateTaskSuite) SetupSuite() {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	s.Require().NoError(err, "could not start postgres container")

	s.T().Cleanup(func() {
		s.Require().NoError(pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	s.Require().NoError(err, "failed to open gorm db")
	s.Require().NoError(db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	s.db = db
	s.service = models.NewTaskService(models.NewTaskStore(db))
}

func (s *UpdateTaskSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY").Error)
}

func (s *UpdateTaskSuite) createTask(userID, title, description string) string {
	task, err := s.service.CreateTask(context.Background(), userID, title, description)
	s.Require().NoError(err, "failed to create initial task")

	return task.ID
}

func (s *UpdateTaskSuite) TestUpdateTaskScenarios() {
	newString := func(value string) *string { return &value }
	newBool := func(value bool) *bool { return &value }

	cases := []struct {
		name            string
		title           *string
		description     *string
		isCompleted     *bool
		wantErr         error
		wantErrContains string
		assertStored    func(t *testing.T, task models.Task)
		assertReturned  func(t *testing.T, task models.TaskModel)
	}{
		{
			name:        "updates all fields successfully",
			title:       newString("Updated title"),
			description: newString("Updated description"),
			isCompleted: newBool(true),
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.True(t, task.IsCompleted)
			},
			assertReturned: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.True(t, task.IsCompleted)
			},
		},
		{
			name:        "updates only title",
			title:       newString("Updated title"),
			description: nil,
			isCompleted: nil,
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
			assertReturned: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Updated title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
		},
		{
			name:        "updates only description",
			title:       nil,
			description: newString("Updated description"),
			isCompleted: nil,
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.False(t, task.IsCompleted)
			},
			assertReturned: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Updated description", task.Description)
				assert.False(t, task.IsCompleted)
			},
		},
		{
			name:        "updates only completion status",
			title:       nil,
			description: nil,
			isCompleted: newBool(true),
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.True(t, task.IsCompleted)
			},
			assertReturned: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.True(t, task.IsCompleted)
			},
		},
		{
			name:        "no-op update leaves values unchanged",
			title:       nil,
			description: nil,
			isCompleted: nil,
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
			assertReturned: func(t *testing.T, task models.TaskModel) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
		},
		{
			name:        "rejects blank title",
			title:       newString("   "),
			description: newString("Updated description"),
			wantErr:     models.ErrInvalidTask,
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
		},
		{
			name:        "rejects blank description",
			title:       newString("Updated title"),
			description: newString("   "),
			wantErr:     models.ErrInvalidTask,
			assertStored: func(t *testing.T, task models.Task) {
				assert.Equal(t, "Original title", task.Title)
				assert.Equal(t, "Original description", task.Description)
				assert.False(t, task.IsCompleted)
			},
		},
		{
			name:            "returns error when task does not exist",
			title:           newString("Updated title"),
			description:     newString("Updated description"),
			wantErrContains: "task not found",
		},
	}

	userStore := models.NewUserStore(s.db)
	userService := models.NewUserService(userStore)

	createdUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	s.Require().NoError(err, "failed to create test user")

	for _, tc := range cases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			var userID string
			var taskID string

			if tc.wantErrContains == "task not found" {
				userID = "missing-user"
				taskID = "missing-task"
			} else {
				userID = createdUser.ID
				taskID = s.createTask(userID, "Original title", "Original description")
			}

			response, err := s.service.UpdateTask(ctx, taskID, userID, tc.title, tc.description, tc.isCompleted)
			if tc.wantErr != nil {
				assert.ErrorIs(s.T(), err, tc.wantErr)
			} else if tc.wantErrContains != "" {
				assert.ErrorContains(s.T(), err, tc.wantErrContains)
			} else {
				assert.NoError(s.T(), err)
			}

			if tc.assertReturned != nil && tc.wantErrContains != "task not found" {
				tc.assertReturned(s.T(), *response)
			}

			if tc.assertStored != nil && tc.wantErrContains != "task not found" {
				var stored models.Task
				s.Require().NoError(s.db.First(&stored, "id = ?", taskID).Error)
				tc.assertStored(s.T(), stored)
			}
		})
	}
}
