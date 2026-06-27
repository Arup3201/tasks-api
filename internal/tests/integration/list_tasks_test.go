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

type ListTasksSuite struct {
	suite.Suite
	db      *gorm.DB
	service *models.TaskService
	userID  string
}

func TestListTasksSuite(t *testing.T) {
	suite.Run(t, new(ListTasksSuite))
}

func (s *ListTasksSuite) SetupSuite() {
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

func (s *ListTasksSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY").Error)

	userStore := models.NewUserStore(s.db)
	userService := models.NewUserService(userStore)
	createdUser, err := userService.CreateUser(context.Background(), "list-user@example.com", "List User", "secret123")
	s.Require().NoError(err, "failed to create test user")
	s.userID = createdUser.ID
}

func (s *ListTasksSuite) createTask(title, description string) string {
	task, err := s.service.CreateTask(context.Background(), s.userID, title, description)
	s.Require().NoError(err, "failed to create test task")

	return task.ID
}

func (s *ListTasksSuite) TestListTasksReturnsUserTasks() {
	firstID := s.createTask("First task", "First description")
	secondID := s.createTask("Second task", "Second description")

	tasks, err := s.service.ListTasks(context.Background(), s.userID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), tasks, 2)

	gotIDs := map[string]models.TaskModel{}
	for _, task := range tasks {
		gotIDs[task.ID] = task
	}

	first, ok := gotIDs[firstID]
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "First task", first.Title)
	assert.Equal(s.T(), "First description", first.Description)
	assert.False(s.T(), first.IsCompleted)

	second, ok := gotIDs[secondID]
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "Second task", second.Title)
	assert.Equal(s.T(), "Second description", second.Description)
}

func (s *ListTasksSuite) TestListTasksReturnsEmptyForUserWithNoTasks() {
	tasks, err := s.service.ListTasks(context.Background(), s.userID)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), tasks)
}

func (s *ListTasksSuite) TestListTasksDoesNotReturnOtherUsersTasks() {
	otherStore := models.NewUserStore(s.db)
	otherService := models.NewUserService(otherStore)
	otherUser, err := otherService.CreateUser(context.Background(), "other-user@example.com", "Other User", "secret123")
	s.Require().NoError(err, "failed to create other user")

	_ = s.createTask("Owner task", "Owner description")
	taskForOther, err := s.service.CreateTask(context.Background(), otherUser.ID, "Other task", "Other description")
	s.Require().NoError(err, "failed to create task for other user")

	tasks, err := s.service.ListTasks(context.Background(), s.userID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), tasks, 1)
	assert.Equal(s.T(), "Owner task", tasks[0].Title)
	assert.NotEqual(s.T(), taskForOther.ID, tasks[0].ID)
}
