package web

import (
	"github.com/gin-gonic/gin"
	"priority-task-manager/manager/internal/services"
	"priority-task-manager/shared/pkg/services/settings"
)

func MakeControllers(
	r *gin.Engine,
	basicAuth settings.BasicAuthSettings,
	services services.Services,
	extracAccountMiddleware gin.HandlerFunc,

	sharedMiddlewares ...gin.HandlerFunc,
) {
	initInternalControllers(
		r,
		basicAuth,
		services,
		sharedMiddlewares...,
	)

	initExternalControllers(
		r,
		services,
		extracAccountMiddleware,
		sharedMiddlewares...,
	)
}

func initInternalControllers(
	r *gin.Engine,
	basicAuth settings.BasicAuthSettings,
	services services.Services,
	sharedMiddlewares ...gin.HandlerFunc,
) {
	alertController := makeAlertController(services.PermissionVersionManger())

	internal := r.Group("/internal", gin.BasicAuth(gin.Accounts{
		basicAuth.User: basicAuth.Password,
	}))

	internal.Use(sharedMiddlewares...)

	internal.POST("/alert", alertController.Handle)
}

func initExternalControllers(
	r *gin.Engine,
	services services.Services,
	extracAccountMiddleware gin.HandlerFunc,
	sharedMiddlewares ...gin.HandlerFunc,
) {
	taskController := makeTaskController(services.Task())

	external := r.Group("/external")
	external.Use(extracAccountMiddleware)
	external.Use(sharedMiddlewares...)

	external.POST("/task", taskController.Handle)
}
