package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Arup3201/gotasks/internal/config"
	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/health"
	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/storages"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/rs/cors"
)

func main() {
	config := config.Load()
	db, err := storages.New(config)
	if err != nil {
		log.Fatal(err)
	}

	userStore := models.NewUserStore(db)
	taskStore := models.NewTaskStore(db)

	userService := models.NewUserService(userStore)
	jwtService := utils.NewJWTService(config)
	taskService := models.NewTaskService(taskStore)

	authController := controllers.NewAuthController(userService, jwtService)
	taskController := controllers.NewTaskController(taskService)
	authMiddleware := middlewares.NewAuthMiddleware(jwtService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", authController.Register)
	mux.HandleFunc("POST /login", authController.Login)

	taskEndpoints := []struct {
		patter string
		fn     http.HandlerFunc
	}{
		{
			patter: "POST /tasks",
			fn:     taskController.CreateTask,
		},
		{
			patter: "PATCH /tasks/{id}",
			fn:     taskController.UpdateTask,
		},
		{
			patter: "GET /tasks",
			fn:     taskController.ListTasks,
		},
		{
			patter: "DELETE /tasks/{id}",
			fn:     taskController.DeleteTask,
		},
	}

	for _, ep := range taskEndpoints {
		mux.Handle(ep.patter, authMiddleware.
			Required(http.
				HandlerFunc(ep.fn),
			))
	}

	healthChecker := health.NewHealthChecker(db)
	mux.HandleFunc("GET /health", healthChecker.HealthHandler)

	server := http.Server{
		Addr: fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: cors.New(cors.Options{
			AllowedMethods: []string{"HEAD", "GET", "POST", "PATCH", "DELETE"},
			AllowedHeaders: []string{"Authorization", "Content-Type"},
		}).Handler(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
