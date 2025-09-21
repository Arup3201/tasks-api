package http

import (
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

func CreateServer(storage storages.TaskRepository) *HttpServer {
	engine := gin.Default()
	engine.Use(middlewares.HttpErrorResponse())
	serviceHandler := task.NewTaskService(storage)
	return &HttpServer{
		engine:       engine,
		routeHandler: getRouteHandler(serviceHandler),
	}
}

func (server *HttpServer) AttachRoutes() {
	server.engine.GET("/tasks", server.routeHandler.GetTasks)
	server.engine.POST("/tasks", server.routeHandler.AddTask)
	server.engine.GET("/tasks/:id", server.routeHandler.GetTask)
	server.engine.PATCH("/tasks/:id", server.routeHandler.UpdateTask)
	server.engine.DELETE("/tasks/:id", server.routeHandler.DeleteTask)
	server.engine.GET("/search/tasks", server.routeHandler.SearchTasks)
}

func (server *HttpServer) Run(host string) error {
	err := server.engine.Run(host + ":" + utils.Config.Port)
	return err
}
