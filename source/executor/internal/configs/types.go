package configs

import "priority-task-manager/shared/pkg/services/settings"

type (
	Settings struct {
		DB       settings.DBSettings       `json:"db"`
		RabbitMQ settings.RabbitMQSettings `json:"rabbitmq"`
	}
)
