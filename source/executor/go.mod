module priority-task-manager/executor

go 1.18

require (
	github.com/jmoiron/sqlx v1.3.1
	github.com/rabbitmq/amqp091-go v1.5.0
	github.com/sirupsen/logrus v1.8.1
	priority-task-manager/shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/lib/pq v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
)

replace priority-task-manager/shared => ../shared
