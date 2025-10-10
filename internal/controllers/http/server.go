package httpController

import (
	"net/http"

	"github.com/Arup3201/gotasks/internal/controllers/http/middlewares"
	"github.com/Arup3201/gotasks/internal/services/domain/task"
	"github.com/Arup3201/gotasks/internal/storages"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	engine       *gin.Engine
	routeHandler *routeHandler
}

var Server = &HttpServer{}

func InitServer(storage storages.TaskRepository) error {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middlewares.HttpErrorResponse())
	engine.Use(middlewares.Authenticate([]string{"/tasks", "/search"}))

	serviceHandler, err := task.NewTaskService(storage)
	if err != nil {
		return err
	}

	Server.engine = engine
	Server.routeHandler = GetRouteHandler(serviceHandler)

	Server.AttachRoutes()

	return nil
}

func (server *HttpServer) AttachRoutes() {
	server.engine.POST("/login", server.routeHandler.Login)
	server.engine.GET("/tasks", server.routeHandler.GetTasks)
	server.engine.POST("/tasks", server.routeHandler.AddTask)
	server.engine.GET("/tasks/:id", server.routeHandler.GetTask)
	server.engine.PATCH("/tasks/:id", server.routeHandler.UpdateTask)
	server.engine.DELETE("/tasks/:id", server.routeHandler.DeleteTask)
	server.engine.GET("/search/tasks", server.routeHandler.SearchTasks)
}

func (server *HttpServer) UpdateLastInsertedId(lastInsertedId int) { // useful for testing setup
	server.routeHandler.serviceHandler.UpdateLastInsertedId(lastInsertedId)
}

func (server *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.engine.ServeHTTP(w, r)
}

func (server *HttpServer) Run(host string) error {
	err := server.engine.Run(host + ":" + utils.Config.Port)
	return err
}
