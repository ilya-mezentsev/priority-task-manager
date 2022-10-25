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

// ProcessTask превичная обработка задачи: вычисляем приоритет на основе роли аккаунта, помещаем в очередь и добавляем в БД
func (s Service) ProcessTask(request Request) error {
	priority, err := s.queuePriority.DetermineMaxPriority(request.Account.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
		}).Errorf("unable to determine max priority: %v", err)

		return err
	}

	t := queue.Task[any]{
		Priority:    priority,
		UUID:        uuid.New().String(),
		AccountHash: request.Account.Hash,
		Type:        request.Task.Type,
		Data:        request.Task.Data,
	}
	err = s.queueImpl.Add(t)
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
			"task":    t,
		}).Errorf("Unable to add task to queue: %v", err)

		return err
	}

	err = s.repository.Add(types.Task{
		UUID:         t.UUID,
		AccountHash:  t.AccountHash,
		Status:       types.InitialStatus,
		AddedToQueue: time.Now(),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"request": request,
			"task":    t,
		}).Errorf("Unable to add task to DB: %v", err)

		return err
	}

	return nil
}
