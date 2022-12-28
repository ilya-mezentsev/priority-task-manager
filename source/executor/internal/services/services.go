package services

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/executor/internal/configs"
	"priority-task-manager/executor/internal/services/consumer/rabbitmq"
	"priority-task-manager/executor/internal/services/worker"
	"priority-task-manager/shared/pkg/repositories"
	sharedRabbitMQ "priority-task-manager/shared/pkg/services/queue/rabbitmq"
)

type Services struct {
	taskExecutor worker.Executor
	workersPool  *worker.Pool
	taskConsumer rabbitmq.Service
}

func MakeServices(settings configs.Settings, db *sqlx.DB) Services {
	s := Services{}

	s.taskExecutor = worker.MakeExecutor(repositories.MakeTaskRepository(db))
	s.workersPool = worker.MakePool(
		settings.WorkersPool.MaxWorkersCount,
		s.taskExecutor,
	)

	channel, err := sharedRabbitMQ.InitChannel(settings.RabbitMQ)
	if err != nil {
		log.Fatalf("Unable to init services. Rabbitmq channel init failed: %v", err)
	}

	// выставляем prefetchCount, чтобы за раз извлекать только одно сообщение из очереди
	err = channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	if err != nil {
		log.Fatalf("Unable to init services. Rabbitmq channel setup failed: %v", err)
	}

	s.taskConsumer = rabbitmq.MakeService(s.workersPool, channel, settings.RabbitMQ)

	return s
}

func (s Services) TaskConsumer() rabbitmq.Service {
	return s.taskConsumer
}
