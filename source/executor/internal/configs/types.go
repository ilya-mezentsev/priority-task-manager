package configs

import "priority-task-manager/shared/pkg/services/settings"

type (
	WorkersPoolSettings struct {
		MaxWorkersCount int `json:"max_workers_count"`
	}

	RedisSettings struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}

	Settings struct {
		DB            settings.DBSettings       `json:"db"`
		RabbitMQ      settings.RabbitMQSettings `json:"rabbitmq"`
		WorkersPool   WorkersPoolSettings       `json:"workers_pool"`
		RedisSettings RedisSettings             `json:"redis"`
	}
)
