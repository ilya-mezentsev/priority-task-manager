package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"priority-task-manager/metrics/internal/configs"
	"priority-task-manager/metrics/internal/services"
	"priority-task-manager/metrics/internal/transport/web"
)

func Main() {
	configsPath := os.Args[1]
	_, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	r := gin.New()
	web.MakeControllers(
		r,
		services.MakeServices(),
	)
}
