package app

import (
	"log"
	"os"
	"priority-task-manager/executor/internal/configs"
	"priority-task-manager/shared/pkg/services/queue/rabbitmq"
)

func Main() {
	configsPath := os.Args[1]
	settings, err := configs.ParseConfigs(configsPath)
	if err != nil {
		log.Fatalf("Unable to parse configs by path %s, got error %v\n", configsPath, err)
	}

	channel, err := rabbitmq.InitChannel(settings.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to init rabbitmq channel: %v", err)
	}

	msgs, err := channel.Consume(
		settings.RabbitMQ.Queue.Name,
		settings.RabbitMQ.Queue.Consumer,
		settings.RabbitMQ.Queue.AutoAck,
		settings.RabbitMQ.Queue.Exclusive,
		settings.RabbitMQ.Queue.NoLocal,
		settings.RabbitMQ.Queue.NoWait,
		nil, // args
	)

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
