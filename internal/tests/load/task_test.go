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
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type LoadSuite struct {
	suite.Suite
	ctx        context.Context
	pg         *testutils.PostgresContainer
	token      string
	httpServer http.Server
}

func TestLoadSuite(t *testing.T) {
	suite.Run(t, new(LoadSuite))
}

func (s *LoadSuite) SetupSuite() {
	var err error

	s.ctx = context.Background()
	s.pg, err = testutils.CreatePostgresContainer(s.ctx)
	s.Require().NoError(err, "could not start postgres container")

	db, err := gorm.Open(postgres.Open(s.pg.ConnectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	s.Require().NoError(err, "failed to open gorm db")
	s.Require().NoError(db.AutoMigrate(&models.User{}, &models.Task{}), "failed to migrate schema")

	s.T().Setenv("HOST", "localhost")
	s.T().Setenv("PORT", "8080")
	s.T().Setenv("JWT_SECRET", "test-secret")
	s.T().Setenv("JWT_ISSUER", "test-issuer")
	config := config.Load()

	userStore := models.NewUserStore(db)
	userService := models.NewUserService(userStore)
	testUser, err := userService.CreateUser(context.Background(), "task-user@example.com", "Task User", "secret123")
	s.Require().NoError(err, "failed to create user")

	taskStore := models.NewTaskStore(db)
	taskService := models.NewTaskService(taskStore)
	jwtService := utils.NewJWTService(config)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	s.token, _ = jwtService.GenerateToken(testUser.ID, testUser.Email)

	mux := http.NewServeMux()
	mux.Handle("POST /tasks", authMiddleware.
		Required(http.
			HandlerFunc(taskController.CreateTask),
		))
	mux.Handle("PATCH /tasks/{id}", authMiddleware.
		Required(http.
			HandlerFunc(taskController.UpdateTask),
		))
	s.httpServer = http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler:      mux,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}
	go s.httpServer.ListenAndServe()
}

func (s *LoadSuite) Teardown() {
	s.Require().NoError(s.pg.Terminate(s.ctx), "could not terminate postgres container")

	timeredCtx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	s.Require().NoError(s.httpServer.Shutdown(timeredCtx))
	defer cancel()
}

func (s *LoadSuite) TestLoadCreateTaskEndpoint() {
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
			req.Header.Set("Authorization", "Bearer "+s.token)

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

	s.T().Logf("Load test results:")
	s.T().Logf("Total requests: %d", requests)
	s.T().Logf("Successful requests: %d", successCount)
	s.T().Logf("Success rate: %.2f%%", float64(successCount)/float64(requests)*100)
	s.T().Logf("Total duration: %v", totalDuration)
	s.T().Logf("Throughput: %.2f requests/second", throughput)
	s.T().Logf("Average latency: %v", avgLatency)
	s.T().Logf("Min latency: %v", minLatency)
	s.T().Logf("Max latency: %v", maxLatency)

	s.Require().True(float64(successCount)/float64(requests) > 0.95, "Success rate should be > 95%")
	s.Require().True(avgLatency < 100*time.Millisecond, "Average latency should be < 100ms")
}

func (s *LoadSuite) TestLoadUpdateTaskEndpoint() {
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

	client := &http.Client{Timeout: 10 * time.Second}

	jsonBody, _ := json.Marshal(controllers.CreateTaskRequest{
		Title:       "Original Title",
		Description: "Original Description",
	})

	req, _ := http.NewRequest("POST", baseUrl+"/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	response, err := client.Do(req)

	s.Require().NoError(err)
	s.Require().Equal(201, response.StatusCode)

	var v controllers.CreateTaskResponse
	json.NewDecoder(response.Body).Decode(&v)

	taskID := v.Task.ID

	newString := func(s string) *string {
		return &s
	}
	newBool := func(bl bool) *bool {
		return &bl
	}

	var wg sync.WaitGroup
	var results = make(chan result, requests)

	worker := func(taskChan <-chan int) {
		defer wg.Done()

		client := &http.Client{Timeout: 10 * time.Second}

		for i := range taskChan {
			startTime := time.Now()

			jsonBody, _ := json.Marshal(controllers.TaskUpdateRequest{
				Title:       newString(fmt.Sprintf("Load Test Task Title %d", i)),
				Description: newString(fmt.Sprintf("Load Test Task Description %d", i)),
				IsCompleted: newBool(true),
			})

			req, _ := http.NewRequest("PATCH", baseUrl+"/tasks/"+taskID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+s.token)

			response, err := client.Do(req)

			res := result{
				duration: time.Since(startTime),
				err:      err,
			}

			if response != nil {
				res.code = response.StatusCode
				response.Body.Close()
			}

			if err != nil {
				s.T().Log(err)
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

	var (
		successCount = 0
		totalLatency time.Duration
		maxLatency   time.Duration
		minLatency   = time.Hour
	)
	for res := range results {
		if res.err == nil && res.code == http.StatusOK {
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

	s.T().Logf("Load test results:")
	s.T().Logf("Total requests: %d", requests)
	s.T().Logf("Successful requests: %d", successCount)
	s.T().Logf("Success rate: %.2f%%", float64(successCount)/float64(requests)*100)
	s.T().Logf("Total duration: %v", totalDuration)
	s.T().Logf("Throughput: %.2f requests/second", throughput)
	s.T().Logf("Average latency: %v", avgLatency)
	s.T().Logf("Min latency: %v", minLatency)
	s.T().Logf("Max latency: %v", maxLatency)

	s.Require().True(float64(successCount)/float64(requests) > 0.95, "Success rate should be > 95%")
	s.Require().True(avgLatency < 100*time.Millisecond, "Average latency should be < 100ms")
}
