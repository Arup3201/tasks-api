package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	controller "github.com/Arup3201/gotasks/internal/controllers/http"
	middleware "github.com/Arup3201/gotasks/internal/controllers/http/middlewares"
	services "github.com/Arup3201/gotasks/internal/services/domain/task"
	"github.com/Arup3201/gotasks/internal/storages"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setUp()
	exitCode := m.Run()
	tearDown()

	os.Exit(exitCode)
}

func setUp() {
	utils.Config.Configure("../../.env.test.local")
}

func tearDown() {
	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("tearDown() failed: %v", err)
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

func router() *gin.Engine {
	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("Storage creation failed: %v", err)
	}

	engine := gin.Default()
	engine.Use(middleware.HttpErrorResponse())
	serviceHandler, err := services.NewTaskService(storage)
	if err != nil {
		log.Fatalf("router() failed: %v", err)
	}

	routeHandler := controller.GetRouteHandler(serviceHandler)

	engine.GET("/tasks", routeHandler.GetTasks)
	engine.POST("/tasks", routeHandler.AddTask)
	engine.GET("/tasks/:id", routeHandler.GetTask)
	engine.PATCH("/tasks/:id", routeHandler.UpdateTask)
	engine.DELETE("/tasks/:id", routeHandler.DeleteTask)
	engine.GET("/search/tasks", routeHandler.SearchTasks)

	return engine
}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}
