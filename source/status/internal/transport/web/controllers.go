package web

import (
	"github.com/gin-gonic/gin"
	"priority-task-manager/status/internal/services"
)

func MakeControllers(
	r *gin.Engine,
	services services.Services,
	sharedMiddlewares ...gin.HandlerFunc,
) {
}
