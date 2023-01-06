package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
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
	uniqueRolesRepository             repositories.NoKeyReader[[]types.Role]
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role]
	avgQueuedWaitingTimeRepository    repositories.Reader[float64, types.Role]
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role]
}

func MakeWaitingTimeService(
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role],
	avgQueuedWaitingTimeRepository repositories.Reader[float64, types.Role],
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role],
) WaitingTimeService {
	return WaitingTimeService{
		avgExtractedWaitingTimeRepository: avgExtractedWaitingTimeRepository,
		avgQueuedWaitingTimeRepository:    avgQueuedWaitingTimeRepository,
		avgCompletedWaitingTimeRepository: avgCompletedWaitingTimeRepository,
	}
}

func (queueWaitingTimeService WaitingTimeService) UpdateExtracted(role types.Role) error {
	extractedTaskWaitingTime, err := queueWaitingTimeService.avgExtractedWaitingTimeRepository.Get(role)
	if err != nil {
		return err
	}

	extractedFromQueueWaitingTime.WithLabelValues(string(role)).Observe(extractedTaskWaitingTime)

	return nil
}

func (queueWaitingTimeService WaitingTimeService) UpdateQueued(role types.Role) error {
	queuedTaskWaitingTime, err := queueWaitingTimeService.avgQueuedWaitingTimeRepository.Get(role)
	if err != nil {
		return err
	}

	queuedWaitingTime.WithLabelValues(string(role)).Observe(queuedTaskWaitingTime)

	return nil
}

func (queueWaitingTimeService WaitingTimeService) UpdateComplete(role types.Role) error {
	completeTaskWaitingTime, err := queueWaitingTimeService.avgCompletedWaitingTimeRepository.Get(role)
	if err != nil {
		return err
	}

	completeWaitingTime.WithLabelValues(string(role)).Observe(completeTaskWaitingTime)

	return nil
}
