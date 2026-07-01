package load

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/Arup3201/gotasks/internal/config"
	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/testutils"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoadCreateTaskEndpoint(t *testing.T) {
	ctx := context.Background()
	pg, err := testutils.CreatePostgresContainer(ctx)
	require.NoError(t, err, "could not start postgres container")

	t.Cleanup(func() {
		require.NoError(t, pg.Terminate(ctx), "could not terminate postgres container")
	})

	db, err := gorm.Open(postgres.Open(pg.ConnectionString), &gorm.Config{})
	require.NoError(t, err, "failed to open gorm db")

	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "8080")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	userStore := models.NewUserStore(db)
	userService := models.NewUserService(userStore)
	testUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	require.NoError(t, err, "failed to create user")

	taskStore := models.NewTaskStore(db)
	taskService := models.NewTaskService(taskStore)
	jwtService := utils.NewJWTService(config)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	accessToken, _ := jwtService.GenerateToken(testUser.ID, testUser.Email)

	mux := http.NewServeMux()
	mux.Handle("POST /tasks", authMiddleware.
		Required(http.
			HandlerFunc(taskController.CreateTask),
		))
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler:      (mux),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}
	go server.ListenAndServe()

	var (
		concurrency = 50
		requests    = 1000
		baseUrl     = "http://localhost:8080"
	)

	type result struct {
		code     int
		duration time.Duration
		err      error
	}

	var wg sync.WaitGroup
	var results = make(chan result, requests)

	worker := func(taskChan <-chan int) {
		defer wg.Done()

		client := &http.Client{Timeout: 10 * time.Second}

		for i := range taskChan {
			startTime := time.Now()

			jsonBody, _ := json.Marshal(controllers.CreateTaskRequest{
				Title:       fmt.Sprintf("Load Test Task Title %d", i),
				Description: fmt.Sprintf("Load Test Task Description %d", i),
			})

			req, _ := http.NewRequest("POST", baseUrl+"/tasks", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+accessToken)

			response, err := client.Do(req)

			res := result{
				duration: time.Since(startTime),
				err:      err,
			}

			if response != nil {
				res.code = response.StatusCode
				response.Body.Close()
			}

			results <- res
		}
	}

	var taskChan = make(chan int, requests)
	for range concurrency {
		wg.Add(1)
		go worker(taskChan)
	}

	startTime := time.Now()
	for i := range requests {
		taskChan <- i
	}
	close(taskChan)

	wg.Wait()
	close(results)

	totalDuration := time.Since(startTime)

	timeredCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	require.NoError(t, server.Shutdown(timeredCtx))
	defer cancel()

	var (
		successCount = 0
		totalLatency time.Duration
		maxLatency   time.Duration
		minLatency   = time.Hour
	)
	for res := range results {
		if res.err == nil && res.code == http.StatusCreated {
			successCount++
		}

		totalLatency += res.duration
		if res.duration > maxLatency {
			maxLatency = res.duration
		}
		if res.duration < minLatency {
			minLatency = res.duration
		}
	}

	avgLatency := totalLatency / time.Duration(requests)
	throughput := float64(successCount) / totalDuration.Seconds()

	t.Logf("Load test results:")
	t.Logf("Total requests: %d", requests)
	t.Logf("Successful requests: %d", successCount)
	t.Logf("Success rate: %.2f%%", float64(successCount)/float64(requests)*100)
	t.Logf("Total duration: %v", totalDuration)
	t.Logf("Throughput: %.2f requests/second", throughput)
	t.Logf("Average latency: %v", avgLatency)
	t.Logf("Min latency: %v", minLatency)
	t.Logf("Max latency: %v", maxLatency)

	// Assertions
	assert.True(t, float64(successCount)/float64(requests) > 0.95, "Success rate should be > 95%")
	assert.True(t, avgLatency < 100*time.Millisecond, "Average latency should be < 100ms")
}
