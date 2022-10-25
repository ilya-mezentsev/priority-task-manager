package web

import (
	"github.com/gin-gonic/gin"
	"priority-task-manager/metrics/internal/services"
)

func MakeControllers(
	r *gin.Engine,
	services services.Services,
	sharedMiddlewares ...gin.HandlerFunc,
) {
}
