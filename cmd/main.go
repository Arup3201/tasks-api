package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Arup3201/gotasks/internal/controllers"
	"github.com/Arup3201/gotasks/internal/health"
	"github.com/Arup3201/gotasks/internal/middlewares"
	"github.com/Arup3201/gotasks/internal/models"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is missing", key)
	}

	return value
}

func PostgresDSN() string {
	host := getEnv("DBHOST")
	port := getEnv("DBPORT")
	user := getEnv("DBUSER")
	password := getEnv("DBPASS")
	db := getEnv("DBNAME")

	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port,
		user, password, db)
}

func main() {
	dsn := PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Task{})

	secretKey := getEnv("JWT_SECRET")
	issuer := getEnv("JWT_ISSUER")

	userStore := models.NewUserStore(db)
	taskStore := models.NewTaskStore(db)

	userService := models.NewUserService(userStore)
	jwtService := utils.NewJWTService(secretKey, issuer)
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

	serverHost := getEnv("HOST")
	serverPort := getEnv("PORT")
	server := http.Server{
		Addr: fmt.Sprintf("%s:%s", serverHost, serverPort),
		Handler: cors.New(cors.Options{
			AllowedMethods: []string{"HEAD", "GET", "POST", "PATCH", "DELETE"},
			AllowedHeaders: []string{"Authorization", "Content-Type"},
		}).Handler(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
