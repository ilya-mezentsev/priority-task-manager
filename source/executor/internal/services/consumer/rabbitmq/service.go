package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/executor/internal/services/worker"
	"priority-task-manager/shared/pkg/services/settings"
	"priority-task-manager/shared/pkg/types"
)

type Service struct {
	workersPool      *worker.Pool
	channel          *amqp.Channel
	settings         settings.RabbitMQSettings
	consumingStopped bool
}

func MakeService(
	workersPool *worker.Pool,
	channel *amqp.Channel,
	settings settings.RabbitMQSettings,
) Service {
	return Service{
		workersPool: workersPool,
		channel:     channel,
		settings:    settings,
	}
}

func (s Service) StartConsume() {
	msgs, err := s.channel.Consume(
		s.settings.Queue.Name,
		s.settings.Queue.Consumer,
		s.settings.Queue.AutoAck,
		s.settings.Queue.Exclusive,
		s.settings.Queue.NoLocal,
		s.settings.Queue.NoWait,
		nil, // args
	)

	if err != nil {
		log.Fatalf("Unable to start consume: %v", err)
	}

	for {
		queuedTask, ok := <-msgs
		if !ok {
			log.Error("Unable to receive next message from rabbitmq channel")
			break
		}

		var task types.Task
		err = json.Unmarshal(queuedTask.Body, &task)
		if err != nil {
			log.WithFields(log.Fields{
				"queuedTask": queuedTask,
			}).Errorf("Unable to decode task: %v", err)
			continue
		}

		log.Infof("Processing task: %s", task.UUID)
		s.workersPool.Exec(task)

		// todo надо ли обрабатывать ошибку?
		_ = queuedTask.Ack(false)
	}
}

func (s Service) StopConsume() {
	err := s.channel.Close()
	if err != nil {
		log.Errorf("Unable to close rabbitmq channel: %v", err)
	}
}
