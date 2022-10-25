package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"priority-task-manager/manager/internal/configs"
	"priority-task-manager/manager/internal/services/permission"
	"priority-task-manager/manager/internal/services/queue/rabbitmq"
	"priority-task-manager/manager/internal/services/queue_priority"
	"priority-task-manager/manager/internal/services/task"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/services/cache"
	sharedRabbitmq "priority-task-manager/shared/pkg/services/queue/rabbitmq"
)

type Services struct {
	permissionResolver      permission.Resolver
	permissionVersionManger *permission.VersionManager
	queuePriority           queue_priority.Service
	queue                   rabbitmq.Service
	task                    task.Service
}

func MakeServices(settings configs.Settings, db *sqlx.DB) Services {
	s := Services{}

	s.permissionResolver = permission.MakeResolver(settings.Permission)
	s.permissionVersionManger = permission.MakeVersionManager(settings.Permission)

	s.queuePriority = queue_priority.MakeService(
		settings.QueuePriorities,
		s.permissionVersionManger,
		cache.MakeService[bool, permission.ResolveRequest](
			s.permissionResolver,
			settings.Permission.CacheLifetimeSeconds,
		),
	)

	channel, err := sharedRabbitmq.InitChannel(settings.RabbitMQ)
	if err != nil {
		panic(fmt.Sprintf("Unable to init services. Initing rabbitmq channel failed"))
	}

	s.queue = rabbitmq.MakeService(
		settings.RabbitMQ.Queue.Name,
		settings.RabbitMQ.PublishingTimeoutSeconds,
		channel,
	)

	s.task = task.MakeService(
		repositories.MakeTaskRepository(db),
		s.queue,
		s.queuePriority,
	)

	return s
}

func (s Services) PermissionVersionManger() *permission.VersionManager {
	return s.permissionVersionManger
}

func (s Services) Task() task.Service {
	return s.task
}
