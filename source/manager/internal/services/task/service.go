package task

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/manager/internal/services/queue"
	"priority-task-manager/manager/internal/services/queue_priority"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"time"
)

type Service struct {
	repository    repositories.Repository[types.Task, string]
	queueImpl     queue.TaskReceiver
	queuePriority queue_priority.Service
}

func MakeService(
	repository repositories.Repository[types.Task, string],
	queueImpl queue.TaskReceiver,
	queuePriority queue_priority.Service,
) Service {
	return Service{
		repository:    repository,
		queueImpl:     queueImpl,
		queuePriority: queuePriority,
	}
}

// ProcessTask первичная обработка задачи: вычисляем приоритет на основе роли аккаунта, помещаем в очередь и добавляем в БД
func (s Service) ProcessTask(request Request) (string, error) {
	priority, err := s.queuePriority.DetermineMaxPriority(request.Account.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
		}).Errorf("unable to determine max priority: %v", err)

		return "", err
	}

	t := types.Task{
		UUID:         uuid.New().String(),
		AccountHash:  request.Account.Hash,
		Type:         request.Task.Type,
		Data:         request.Task.Data,
		Status:       types.InitialStatus,
		AddedToQueue: time.Now(),
	}
	err = s.queueImpl.Add(t, priority)
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
			"task":    t,
		}).Errorf("Unable to add task to queue: %v", err)

		return "", err
	}

	err = s.repository.Add(t)
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
			"task":    t,
		}).Errorf("Unable to add task to DB: %v", err)

		return "", err
	}

	return t.UUID, nil
}
