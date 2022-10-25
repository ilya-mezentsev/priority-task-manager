package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"priority-task-manager/manager/internal/configs"
	"priority-task-manager/manager/internal/services"
	"priority-task-manager/manager/internal/transport/web"
	"priority-task-manager/shared/pkg/middlewares"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/services/account"
	"priority-task-manager/shared/pkg/services/db/connection"
	myLogger "priority-task-manager/shared/pkg/services/log"
)

func init() {
	myLogger.Configure()
}

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	db := connection.MustGetConnection(settings.DB)

	sharedMiddlewares := []gin.HandlerFunc{
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: myLogger.Formatter,
			Output:    myLogger.Writer{},
		}),

		gin.Recovery(),
	}

	r := gin.New()
	web.MakeControllers(
		r,
		settings.BasicAuth,

		services.MakeServices(settings, db),

		middlewares.MakeExtractAccount(
			account.MakeService(
				repositories.MakeAccountRepository(db),
			),
		).Extract(),

		sharedMiddlewares...,
	)

	err = r.Run(fmt.Sprintf("%s:%d", settings.Web.Domain, settings.Web.Port))
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
