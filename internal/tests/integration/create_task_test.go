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

type CreateTaskSuite struct {
	suite.Suite
	db      *gorm.DB
	service *models.TaskService
	userID  string
}

func TestCreateTaskSuite(t *testing.T) {
	suite.Run(t, new(CreateTaskSuite))
}

func (s *CreateTaskSuite) SetupSuite() {
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

func (s *CreateTaskSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY").Error)

	userStore := models.NewUserStore(s.db)
	userService := models.NewUserService(userStore)
	createdUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	s.Require().NoError(err, "failed to create test user")
	s.userID = createdUser.ID
}

func (s *CreateTaskSuite) TestCreateTaskSuccess() {
	task, err := s.service.CreateTask(context.Background(), s.userID, "Write tests", "Add unit tests")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), task)
	assert.Equal(s.T(), s.userID, task.UserID)
	assert.Equal(s.T(), "Write tests", task.Title)
	assert.Equal(s.T(), "Add unit tests", task.Description)
	assert.NotEmpty(s.T(), task.ID)
	assert.False(s.T(), task.CreatedAt.IsZero())
	assert.False(s.T(), task.UpdatedAt.IsZero())

	var count int64
	assert.NoError(s.T(), s.db.Model(&models.Task{}).Where("id = ?", task.ID).Count(&count).Error)
	assert.Equal(s.T(), int64(1), count)
}

func (s *CreateTaskSuite) TestCreateTaskInvalidTitle() {
	task, err := s.service.CreateTask(context.Background(), s.userID, "   ", "Valid description")
	assert.ErrorIs(s.T(), err, models.ErrInvalidTask)
	assert.Nil(s.T(), task)

	var count int64
	assert.NoError(s.T(), s.db.Model(&models.Task{}).Count(&count).Error)
	assert.Equal(s.T(), int64(0), count)
}

func (s *CreateTaskSuite) TestCreateTaskInvalidDescription() {
	task, err := s.service.CreateTask(context.Background(), s.userID, "Valid title", "   ")
	assert.ErrorIs(s.T(), err, models.ErrInvalidTask)
	assert.Nil(s.T(), task)

	var count int64
	assert.NoError(s.T(), s.db.Model(&models.Task{}).Count(&count).Error)
	assert.Equal(s.T(), int64(0), count)
}

func (s *CreateTaskSuite) TestCreateTaskWhitespaceValuesAreAccepted() {
	title := "  Write tests  "
	description := "  Add unit tests  "

	task, err := s.service.CreateTask(context.Background(), s.userID, title, description)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), task)
	assert.Equal(s.T(), title, task.Title)
	assert.Equal(s.T(), description, task.Description)

	var stored models.Task
	assert.NoError(s.T(), s.db.First(&stored, "id = ?", task.ID).Error)
	assert.Equal(s.T(), title, stored.Title)
	assert.Equal(s.T(), description, stored.Description)
}
