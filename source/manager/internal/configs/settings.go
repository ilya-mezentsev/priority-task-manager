package configs

import (
	"priority-task-manager/manager/internal/services/permission"
	"priority-task-manager/manager/internal/services/queue_priority"
	"priority-task-manager/shared/pkg/services/settings"
)

type Settings struct {
	Web      settings.WebSettings      `json:"web"`
	DB       settings.DBSettings       `json:"db"`
	RabbitMQ settings.RabbitMQSettings `json:"rabbitmq"`

	Permission      permission.Settings        `json:"permission"`
	QueuePriorities []queue_priority.Settings  `json:"queue_priorities"`
	BasicAuth       settings.BasicAuthSettings `json:"basic_auth"`
}
