package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type createTaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type CreateTaskTestSuite struct {
	suite.Suite
	controller *controllers.TaskController
	db         *gorm.DB
	userSvc    *models.UserService
	taskSvc    *models.TaskService
	userID     string
}

func TestCreateTaskSuite(t *testing.T) {
	suite.Run(t, new(CreateTaskTestSuite))
}

func (s *CreateTaskTestSuite) SetupSuite() {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		s.Require().NoError(pg.Terminate(ctx))
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	s.Require().NoError(err)

	s.Require().NoError(db.AutoMigrate(&models.User{}, &models.Task{}))

	s.db = db
	s.userSvc = models.NewUserService(models.NewUserStore(db))
	s.taskSvc = models.NewTaskService(models.NewTaskStore(db))
	s.controller = controllers.NewTaskController(s.taskSvc)
}

func (s *CreateTaskTestSuite) SetupTest() {
	s.Require().NoError(s.db.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE").Error)

	user, err := s.userSvc.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	s.Require().NoError(err)
	s.userID = user.ID
}

type errorTaskStore struct{}

func (e *errorTaskStore) Create(ctx context.Context, id, userID, title, description string, isCompleted bool) error {
	return errors.New("store failure")
}

func (e *errorTaskStore) Get(ctx context.Context, id, userID string) (*models.TaskModel, error) {
	return nil, errors.New("store failure")
}

func (s *CreateTaskTestSuite) TestCreateTaskEndpoint() {
	cases := []struct {
		name         string
		body         string
		withAuth     bool
		useFaultySvc bool
		wantStatus   int
		wantError    string
		wantTitle    string
		wantDesc     string
		wantComplete bool
	}{
		{
			name:         "success",
			body:         `{"title":"Write tests","description":"Add unit tests"}`,
			withAuth:     true,
			wantStatus:   http.StatusCreated,
			wantTitle:    "Write tests",
			wantDesc:     "Add unit tests",
			wantComplete: false,
		},
		{
			name:       "missing auth",
			body:       `{"title":"Write tests","description":"Add unit tests"}`,
			withAuth:   false,
			wantStatus: http.StatusUnauthorized,
			wantError:  "not authenticated",
		},
		{
			name:       "invalid json",
			body:       `{"title":"Write tests","description":"Add unit tests"`,
			withAuth:   true,
			wantStatus: http.StatusBadRequest,
			wantError:  "json parse error",
		},
		{
			name:       "invalid title",
			body:       `{"title":"   ","description":"Valid description"}`,
			withAuth:   true,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid task title or description provided",
		},
		{
			name:       "invalid description",
			body:       `{"title":"Valid title","description":"   "}`,
			withAuth:   true,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid task title or description provided",
		},
		{
			name:         "server error",
			body:         `{"title":"Write tests","description":"Add unit tests"}`,
			withAuth:     true,
			useFaultySvc: true,
			wantStatus:   http.StatusInternalServerError,
			wantError:    "server error",
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tc.body))
			if tc.withAuth {
				req = req.WithContext(context.WithValue(req.Context(), "user_id", s.userID))
			}

			rec := httptest.NewRecorder()

			controller := s.controller
			if tc.useFaultySvc {
				controller = controllers.NewTaskController(models.NewTaskService(&errorTaskStore{}))
			}

			controller.CreateTask(rec, req)

			assert.Equal(s.T(), tc.wantStatus, rec.Code)

			if tc.wantError != "" {
				assert.Contains(s.T(), rec.Body.String(), tc.wantError)
				return
			}

			var resp createTaskResponse
			require.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
			assert.Equal(s.T(), tc.wantTitle, resp.Title)
			assert.Equal(s.T(), tc.wantDesc, resp.Description)
			assert.Equal(s.T(), tc.wantComplete, resp.IsCompleted)
			assert.NotEmpty(s.T(), resp.ID)
		})
	}
}
