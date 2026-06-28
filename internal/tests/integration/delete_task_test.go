package integration

import (
	"context"
	"testing"

	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DeleteTaskSuite struct {
	suite.Suite
	db      *gorm.DB
	service *models.TaskService
}

func TestDeleteTaskSuite(t *testing.T) {
	suite.Run(t, new(DeleteTaskSuite))
}

func (s *DeleteTaskSuite) SetupSuite() {
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

func (s *DeleteTaskSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY").Error)
}

func (s *DeleteTaskSuite) createUser(email, name, password string) string {
	userStore := models.NewUserStore(s.db)
	userService := models.NewUserService(userStore)

	user, err := userService.CreateUser(context.Background(), email, name, password)
	s.Require().NoError(err, "failed to create test user")

	return user.ID
}

func (s *DeleteTaskSuite) createTask(userID, title, description string) string {
	task, err := s.service.CreateTask(context.Background(), userID, title, description)
	s.Require().NoError(err, "failed to create initial task")

	return task.ID
}

func (s *DeleteTaskSuite) TestDeleteTaskScenarios() {
	cases := []struct {
		name                     string
		setup                    func() (taskID string, userID string)
		wantErr                  error
		wantErrContains          string
		expectCount              int64
		expectTaskPresent        bool
		expectRemainingTaskTitle string
	}{
		{
			name: "deletes task successfully",
			setup: func() (string, string) {
				userID := s.createUser("delete-owner@example.com", "Delete Owner", "secret123")
				taskID := s.createTask(userID, "Task to delete", "Delete this task")
				return taskID, userID
			},
			wantErr:           nil,
			expectCount:       0,
			expectTaskPresent: false,
		},
		{
			name: "returns not found when task does not exist",
			setup: func() (string, string) {
				userID := s.createUser("missing-task-user@example.com", "Missing Task User", "secret123")
				return "missing-task-id", userID
			},
			wantErr:           models.ErrTaskNotFound,
			expectCount:       0,
			expectTaskPresent: false,
		},
		{
			name: "returns not found when user does not own task",
			setup: func() (string, string) {
				ownerID := s.createUser("owner@example.com", "Task Owner", "secret123")
				otherID := s.createUser("wrong-user@example.com", "Wrong User", "secret123")
				taskID := s.createTask(ownerID, "Owner task", "Owned by first user")
				return taskID, otherID
			},
			wantErr:                  models.ErrTaskNotFound,
			expectCount:              1,
			expectTaskPresent:        true,
			expectRemainingTaskTitle: "Owner task",
		},
		{
			name: "deletes only the targeted task and leaves other tasks intact",
			setup: func() (string, string) {
				userID := s.createUser("multi-task-user@example.com", "Multi Task User", "secret123")
				taskID1 := s.createTask(userID, "First task", "Keep second task")
				s.createTask(userID, "Second task", "Should remain")
				return taskID1, userID
			},
			wantErr:                  nil,
			expectCount:              1,
			expectTaskPresent:        false,
			expectRemainingTaskTitle: "Second task",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			taskID, userID := tc.setup()

			err := s.service.DeleteTask(ctx, taskID, userID)
			if tc.wantErr != nil {
				assert.ErrorIs(s.T(), err, tc.wantErr)
			} else if tc.wantErrContains != "" {
				assert.ErrorContains(s.T(), err, tc.wantErrContains)
			} else {
				assert.NoError(s.T(), err)
			}

			var count int64
			s.Require().NoError(s.db.Model(&models.Task{}).Count(&count).Error)
			assert.Equal(s.T(), tc.expectCount, count)

			storedErr := s.db.First(&models.Task{}, "id = ?", taskID).Error
			if tc.expectTaskPresent {
				assert.NoError(s.T(), storedErr)
			} else {
				assert.ErrorIs(s.T(), storedErr, gorm.ErrRecordNotFound)
			}

			if tc.expectRemainingTaskTitle != "" {
				var remaining models.Task
				assert.NoError(s.T(), s.db.First(&remaining, "title = ?", tc.expectRemainingTaskTitle).Error)
				assert.Equal(s.T(), tc.expectRemainingTaskTitle, remaining.Title)
			}

			if count > 0 {
				s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks").Error)
			}
		})
	}
}
