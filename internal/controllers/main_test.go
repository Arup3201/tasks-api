package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	httpController "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/Arup3201/gotasks/internal/entities/task"
	"github.com/Arup3201/gotasks/internal/storages"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	tearDown := setUp()

	exitCode := m.Run()
	tearDown()

	os.Exit(exitCode)
}

func setUp() func() {
	utils.Config.Configure("../../.env.test.local")

	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("storage create error: %v", err)
	}

	controllers.InitServer(storage)

	return func() {
		cleanDB()
		storage.Close()
	}
}

func prepareDBTasks(n int) {
	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("storage create error: %v", err)
	}

	tasks := generateTasks(n)
	for _, task := range tasks {
		storage.Insert(task.Id, task.Title, task.Description)
	}
	httpController.Server.UpdateLastInsertedId(n)
}

func cleanDB() {
	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("storage create error: %v", err)
	}

	tasks, err := storage.List()
	if err != nil {
		log.Fatalf("tearDown() failed: %v", err)
	}

	for _, task := range tasks {
		_, err = storage.Delete(task.Id)
		if err != nil {
			log.Fatalf("tearDown() failed: %v", err)
		}
	}
}

func generateTasks(n int) []task.Task {
	tasks := []task.Task{}
	for i := range n {
		task := task.Task{
			Id:          i + 1,
			Title:       fmt.Sprintf("title - %d", rand.Intn(9999)),
			Description: fmt.Sprintf("description - %d", rand.Intn(9999)),
			IsCompleted: false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	controllers.Server.ServeHTTP(writer, request)
	return writer
}
