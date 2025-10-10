package httpController

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	httperrors "github.com/Arup3201/gotasks/internal/controllers/http/errors"
	"github.com/Arup3201/gotasks/internal/errors"
	"github.com/Arup3201/gotasks/internal/services"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
)

type CreateTask struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type routeHandler struct {
	serviceHandler services.ServiceHandler
}

func GetRouteHandler(handler services.ServiceHandler) *routeHandler {
	return &routeHandler{
		serviceHandler: handler,
	}
}

func (handler *routeHandler) Login(c *gin.Context) {
	var credential struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credential); err != nil {
		c.Error(httperrors.InternalServerError(fmt.Errorf("c.BindJSON failed with error %v", err)))
		return
	}

	client := gocloak.NewClient(utils.Config.KeycloakServerUrl)
	ctx := context.Background()
	token, err := client.Login(ctx, utils.Config.KeycloakClientId, utils.Config.KeycloakClientSecret, utils.Config.KeycloakRealName,
		credential.Username, credential.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		c.Error(httperrors.IncorrectCredentialError())
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
	})
}

func (handler *routeHandler) GetTasks(c *gin.Context) {
	tasks, err := handler.serviceHandler.GetAllTasks()
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (handler *routeHandler) AddTask(c *gin.Context) {
	var payload CreateTask

	if err := c.BindJSON(&payload); err != nil {
		c.Error(httperrors.InternalServerError(fmt.Errorf("c.BindJSON failed with error %v", err)))
		return
	}

	if payload.Title == nil {
		c.Error(httperrors.MissingBodyError(httperrors.ErrorField{
			Field:  "title",
			Reason: "Task 'title' is required",
		}))
		return
	}
	if payload.Description == nil {
		c.Error(httperrors.MissingBodyError(httperrors.ErrorField{
			Field:  "description",
			Reason: "Task 'description' is required",
		}))
		return
	}

	newTask, err := handler.serviceHandler.CreateTask(*payload.Title, *payload.Description)
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}

	c.IndentedJSON(http.StatusCreated, newTask)
}

func (handler *routeHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(httperrors.InvalidRequestParamError(httperrors.ErrorField{
			Field:  "id",
			Reason: "'id' should be a valid integer",
		}))
		return
	}

	task, err := handler.serviceHandler.GetTask(id)
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func (handler *routeHandler) UpdateTask(c *gin.Context) {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(httperrors.InvalidRequestParamError(httperrors.ErrorField{
			Field:  "id",
			Reason: "'id' should be a valid integer",
		}))
		return
	}
	var payload services.UpdateTaskData
	if err := c.BindJSON(&payload); err != nil {
		c.Error(httperrors.InternalServerError(fmt.Errorf("c.BindJSON failed with error %v", err)))
		return
	}

	if payload.Title == nil && payload.Description == nil && payload.IsCompleted == nil {
		c.Error(httperrors.NoOpError())
		return
	}

	editedTask, err := handler.serviceHandler.UpdateTask(id, payload)
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, editedTask)
}

func (handler *routeHandler) DeleteTask(c *gin.Context) {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(httperrors.InvalidRequestParamError(httperrors.ErrorField{
			Field:  "id",
			Reason: "'id' should be a valid integer",
		}))
		return
	}

	taskId, err := handler.serviceHandler.DeleteTask(id)
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, taskId)
}

func (handler *routeHandler) SearchTasks(c *gin.Context) {
	var query string = c.Query("q")
	if query == "" {
		c.Error(httperrors.InvalidRequestParamError(httperrors.ErrorField{
			Field:  "q",
			Reason: "query param 'q' is required",
		}))
		return
	}

	tasks, err := handler.serviceHandler.SearchTasks(query)
	if err != nil {
		appError, ok := err.(*errors.AppError)
		if ok {
			c.Error(httperrors.FromAppError(appError))
		} else {
			c.Error(httperrors.InternalServerError(err))
		}
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}
