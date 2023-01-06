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
	extractedFromQueueWaitingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "extracted_tasks_waiting_time",
	}, []string{"role"})

	queuedWaitingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "queued_tasks_waiting_time",
	}, []string{"role"})

	completeWaitingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "complete_waiting_time",
	}, []string{"role"})
)

type WaitingTimeService struct {
	uniqueRoles                       []types.Role
	uniqueRolesRepository             repositories.NoKeyReader[[]types.Role]
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role]
	avgQueuedWaitingTimeRepository    repositories.Reader[float64, types.Role]
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role]
}

func MustMakeWaitingTimeService(
	uniqueRolesRepository repositories.NoKeyReader[[]types.Role],
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role],
	avgQueuedWaitingTimeRepository repositories.Reader[float64, types.Role],
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role],
) WaitingTimeService {
	waitingTimeService := WaitingTimeService{
		uniqueRolesRepository:             uniqueRolesRepository,
		avgExtractedWaitingTimeRepository: avgExtractedWaitingTimeRepository,
		avgQueuedWaitingTimeRepository:    avgQueuedWaitingTimeRepository,
		avgCompletedWaitingTimeRepository: avgCompletedWaitingTimeRepository,
	}

	uniqueRoles, err := waitingTimeService.uniqueRolesRepository.Get()
	if err != nil {
		log.Fatalf("Unable to get unique roles: %v", err)
	}

	waitingTimeService.uniqueRoles = uniqueRoles

	return waitingTimeService
}

func (queueWaitingTimeService WaitingTimeService) UpdateForEachRole() error {
	var (
		wg       sync.WaitGroup
		updaters = []func(types.Role, chan<- error){
			queueWaitingTimeService.updateExtractedForRole,
			queueWaitingTimeService.updateQueuedForRole,
			queueWaitingTimeService.updateCompleteForRole,
		}
		errors = make(chan error, len(queueWaitingTimeService.uniqueRoles)*len(updaters))
	)

	for _, role := range queueWaitingTimeService.uniqueRoles {
		for _, updater := range updaters {
			wg.Add(1)
			go func(role types.Role, updater func(types.Role, chan<- error)) {
				defer wg.Done()
				updater(role, errors)
			}(role, updater)
		}
	}

	wg.Wait()

	if len(errors) > 0 {
		// возвращаем просто первую ошибку
		return <-errors
	}

	return nil
}

func (queueWaitingTimeService WaitingTimeService) updateExtractedForRole(
	role types.Role,
	errors chan<- error,
) {
	extractedTaskWaitingTime, err := queueWaitingTimeService.avgExtractedWaitingTimeRepository.Get(role)
	if err != nil {
		errors <- err
		return
	}

	extractedFromQueueWaitingTime.WithLabelValues(string(role)).Observe(extractedTaskWaitingTime)
}

func (queueWaitingTimeService WaitingTimeService) updateQueuedForRole(
	role types.Role,
	errors chan<- error,
) {
	queuedTaskWaitingTime, err := queueWaitingTimeService.avgQueuedWaitingTimeRepository.Get(role)
	if err != nil {
		errors <- err
		return
	}

	queuedWaitingTime.WithLabelValues(string(role)).Observe(queuedTaskWaitingTime)
}

func (queueWaitingTimeService WaitingTimeService) updateCompleteForRole(
	role types.Role,
	errors chan<- error,
) {
	completeTaskWaitingTime, err := queueWaitingTimeService.avgCompletedWaitingTimeRepository.Get(role)
	if err != nil {
		errors <- err
		return
	}

	completeWaitingTime.WithLabelValues(string(role)).Observe(completeTaskWaitingTime)
}
