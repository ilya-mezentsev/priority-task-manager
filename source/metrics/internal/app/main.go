package app

import (
	"log"
	"os"
	"priority-task-manager/metrics/internal/configs"
	"priority-task-manager/metrics/internal/services"
	"priority-task-manager/metrics/internal/transport/web"
	"priority-task-manager/shared/pkg/services/db/connection"
)

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	ss := services.MakeServices(connection.MustGetConnection(settings.DB))
	go ss.StartObserveMetrics()

	web.MakeControllers(settings.Web)
}
