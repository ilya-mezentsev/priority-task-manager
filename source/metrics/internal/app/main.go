package app

import (
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"priority-task-manager/metrics/internal/configs"
	"priority-task-manager/metrics/internal/services"
	"priority-task-manager/metrics/internal/transport/web"
	"priority-task-manager/shared/pkg/services/db/connection"
	myLogger "priority-task-manager/shared/pkg/services/log"
	sharedSettings "priority-task-manager/shared/pkg/services/settings"
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

	ss := services.MakeServices(
		connection.MustGetConnection(settings.DB),
		redisClient(settings.Redis),
	)

	go ss.StartObserveMetrics()

	web.MakeControllers(settings.Web)
}

func redisClient(redisSettings sharedSettings.RedisSettings) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisSettings.Address,
		Password: redisSettings.Password,
		DB:       redisSettings.DB,
	})
}
