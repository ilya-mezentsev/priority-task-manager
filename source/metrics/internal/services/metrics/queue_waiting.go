package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"sync"
)

var (
	inQueueWaitingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tasks_in_queue_waiting_time",
	}, []string{"role"})
)

type QueueWaitingTimeService struct {
	uniqueRoles                   []types.Role
	uniqueRolesRepository         repositories.NoKeyReader[[]types.Role]
	avgQueueWaitingTimeRepository repositories.Reader[float64, types.Role]
}

func MustMakeQueueWaitingTimeService(
	uniqueRolesRepository repositories.NoKeyReader[[]types.Role],
	avgQueueWaitingTimeRepository repositories.Reader[float64, types.Role],
) QueueWaitingTimeService {
	queueWaitingTimeService := QueueWaitingTimeService{
		uniqueRolesRepository:         uniqueRolesRepository,
		avgQueueWaitingTimeRepository: avgQueueWaitingTimeRepository,
	}

	uniqueRoles, err := queueWaitingTimeService.uniqueRolesRepository.Get()
	if err != nil {
		log.Fatalf("Unable to get unique roles: %v", err)
	}

	queueWaitingTimeService.uniqueRoles = uniqueRoles

	return queueWaitingTimeService
}

func (queueWaitingTimeService QueueWaitingTimeService) UpdateForEachRole() error {
	var (
		wg     sync.WaitGroup
		errors = make(chan error, len(queueWaitingTimeService.uniqueRoles))
	)
	for _, role := range queueWaitingTimeService.uniqueRoles {
		wg.Add(1)
		go func(role types.Role) {
			defer wg.Done()
			queueWaitingTimeService.updateForRole(role, errors)
		}(role)
	}

	wg.Wait()

	if len(errors) > 0 {
		// возвращаем просто первую ошибку
		return <-errors
	}

	return nil
}

func (queueWaitingTimeService QueueWaitingTimeService) updateForRole(
	role types.Role,
	errors chan<- error,
) {
	waitingTime, err := queueWaitingTimeService.avgQueueWaitingTimeRepository.Get(role)
	if err != nil {
		errors <- err
	}

	inQueueWaitingTime.WithLabelValues(string(role)).Observe(waitingTime)
}
