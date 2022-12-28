package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"priority-task-manager/shared/pkg/types"
	"time"
)

type Service struct {
	queueName                string
	publishingTimeoutSeconds int
	channel                  *amqp.Channel
}

func MakeService(queueName string, publishingTimeoutSeconds int, channel *amqp.Channel) Service {
	return Service{
		queueName:                queueName,
		publishingTimeoutSeconds: publishingTimeoutSeconds,
		channel:                  channel,
	}
}

func (s Service) Add(task types.Task, priority int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.publishingTimeoutSeconds)*time.Second)
	defer cancel()

	encodedTask, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.channel.PublishWithContext(
		ctx,
		"", // exchange
		s.queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        encodedTask,
			Priority:    uint8(priority),
		},
	)
}
