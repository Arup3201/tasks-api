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
	engine := gin.Default()
	engine.Use(middlewares.HttpErrorResponse())

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
	server.engine.GET("/tasks", server.routeHandler.GetTasks)
	server.engine.POST("/tasks", server.routeHandler.AddTask)
	server.engine.GET("/tasks/:id", server.routeHandler.GetTask)
	server.engine.PATCH("/tasks/:id", server.routeHandler.UpdateTask)
	server.engine.DELETE("/tasks/:id", server.routeHandler.DeleteTask)
	server.engine.GET("/search/tasks", server.routeHandler.SearchTasks)
}

func (server *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.engine.ServeHTTP(w, r)
}

func (server *HttpServer) Run(host string) error {
	err := server.engine.Run(host + ":" + utils.Config.Port)
	return err
}
