package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/shared/pkg/services/settings"
)

// InitChannel для подключения к rabbitmq и создания очереди
func InitChannel(settings settings.RabbitMQSettings) (*amqp.Channel, error) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d/", settings.User, settings.Password, settings.Host, settings.Port))
	if err != nil {
		log.WithFields(log.Fields{
			"settings": settings,
		}).Errorf("Unable to create connection to rabbitmq server: %v", err)

		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Unable to create channel: %v", err)

		return nil, err
	}

	_, err = ch.QueueDeclare(
		settings.Queue.Name,
		settings.Queue.Durable,
		settings.Queue.AutoDelete,
		settings.Queue.Exclusive,
		settings.Queue.NoWait,
		amqp.Table{
			"x-max-priority": settings.Queue.MaxPriority,
		}, // arguments
	)
	if err != nil {
		log.WithFields(log.Fields{
			"queue": settings.Queue,
		}).Errorf("Unable to declare queue: %v", err)

		return nil, err
	}

	return ch, nil
}
