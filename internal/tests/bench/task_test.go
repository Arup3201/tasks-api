package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arup3201/gotasks/internal/config"
	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func BenchmarkCreateTaskEndpoint(b *testing.B) {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	require.NoError(b, err, "could not start postgres container")

	b.Cleanup(func() {
		require.NoError(b, pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	require.NoError(b, err, "failed to open gorm db")

	require.NoError(b, db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	userStore := models.NewUserStore(db)
	userService := models.NewUserService(userStore)
	testUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	require.NoError(b, err, "failed to create user")

	b.Setenv("JWT_SECRET", "test-secret")
	b.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	taskStore := models.NewTaskStore(db)
	taskService := models.NewTaskService(taskStore)
	jwtService := utils.NewJWTService(config)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	jsonData, _ := json.Marshal(controllers.CreateTaskRequest{
		Title:       "Benchmark Task",
		Description: "Detail of the task",
	})

	b.ResetTimer()
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.Handle("POST /tasks", authMiddleware.Required(http.HandlerFunc(taskController.CreateTask)))

	accessToken, _ := jwtService.GenerateToken(testUser.ID, testUser.Email)

	for b.Loop() {
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		require.Equal(b, 201, rec.Result().StatusCode)
	}
}

func BenchmarkUpdateTaskEndpoint(b *testing.B) {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	require.NoError(b, err, "could not start postgres container")

	b.Cleanup(func() {
		require.NoError(b, pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	require.NoError(b, err, "failed to open gorm db")

	require.NoError(b, db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	userStore := models.NewUserStore(db)
	userService := models.NewUserService(userStore)
	testUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	require.NoError(b, err, "failed to create user")

	b.Setenv("JWT_SECRET", "test-secret")
	b.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	taskStore := models.NewTaskStore(db)
	taskService := models.NewTaskService(taskStore)
	jwtService := utils.NewJWTService(config)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	task, _ := taskService.CreateTask(ctx, testUser.ID, "Benchmark Task", "Benchmark task description")

	b.ResetTimer()
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.Handle("PATCH /tasks/{id}",
		authMiddleware.Required(http.HandlerFunc(taskController.UpdateTask)))

	accessToken, _ := jwtService.GenerateToken(testUser.ID, testUser.Email)

	newString := func(s string) *string {
		return &s
	}
	newBool := func(bl bool) *bool {
		return &bl
	}

	jsonData, _ := json.Marshal(controllers.TaskUpdateRequest{
		Title:       newString("Benchmark Task"),
		Description: newString("Detail of the task"),
		IsCompleted: newBool(true),
	})

	for b.Loop() {
		req, _ := http.NewRequest("PATCH", "/tasks/"+task.ID, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		require.Equal(b, 200, rec.Result().StatusCode)
	}
}

func BenchmarkListTasksEndpoint(b *testing.B) {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	require.NoError(b, err, "could not start postgres container")

	b.Cleanup(func() {
		require.NoError(b, pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	require.NoError(b, err, "failed to open gorm db")

	require.NoError(b, db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	userStore := models.NewUserStore(db)
	userService := models.NewUserService(userStore)
	testUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	require.NoError(b, err, "failed to create user")

	b.Setenv("JWT_SECRET", "test-secret")
	b.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	taskStore := models.NewTaskStore(db)
	taskService := models.NewTaskService(taskStore)
	jwtService := utils.NewJWTService(config)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	for i := 0; i < 100000; i++ {
		taskService.CreateTask(ctx, testUser.ID, "Benchmark Task", "Benchmark task description")
	}

	b.ResetTimer()
	b.ReportAllocs()

	mux := http.NewServeMux()
	mux.Handle("GET /tasks",
		authMiddleware.Required(http.HandlerFunc(taskController.ListTasks)))

	accessToken, _ := jwtService.GenerateToken(testUser.ID, testUser.Email)

	for b.Loop() {
		req, _ := http.NewRequest("GET", "/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		require.Equal(b, 200, rec.Result().StatusCode)
	}
}
