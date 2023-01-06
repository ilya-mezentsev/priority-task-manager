package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"priority-task-manager/manager/internal/services/task"
	"priority-task-manager/shared/pkg/types"
)

type TaskController struct {
	service task.Service
}

func makeTaskController(service task.Service) TaskController {
	return TaskController{service: service}
}

func (tc TaskController) Handle(context *gin.Context) {
	var r task.Payload
	if err := context.ShouldBindJSON(&r); err != nil {
		context.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	account, found := context.Get("account")
	if !found {
		// такого быть не должно
		context.String(http.StatusInternalServerError, "Server got itself in trouble")
		return
	}

	taskUUID, err := tc.service.ProcessTask(task.Request{
		Account: account.(types.Account),
		Task:    r,
	})
	if err != nil {
		context.String(http.StatusInternalServerError, "Server got itself in trouble")
	} else {
		context.Header("X-Created-Task-UUID", taskUUID)
		context.String(http.StatusCreated, "OK")
	}
}
