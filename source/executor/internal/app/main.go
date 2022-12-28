package app

import (
	"os"
	"priority-task-manager/executor/internal/configs"
	"priority-task-manager/executor/internal/services"
	"priority-task-manager/shared/pkg/services/db/connection"

	log "github.com/sirupsen/logrus"
)

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	db := connection.MustGetConnection(settings.DB)
	ss := services.MakeServices(settings, db)

	log.Info("Starting consume tasks from queue")
	ss.TaskConsumer().StartConsume()
}
