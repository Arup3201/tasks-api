package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	controllers "github.com/Arup3201/gotasks/internal/controllers/http"
	"github.com/Arup3201/gotasks/internal/storages"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	tearDown := setUp()
	defer tearDown()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func setUp() func() {
	utils.Config.Configure("../../.env.test.local")

	storage, err := storages.New(storages.Postgres)
	if err != nil {
		log.Fatalf("tearDown() failed: %v", err)
	}

	controllers.InitServer(storage)

	return func() {
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

		storage.Close()
	}
}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	writer := httptest.NewRecorder()
	controllers.Server.ServeHTTP(writer, request)
	return writer
}
