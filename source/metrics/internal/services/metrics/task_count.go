package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"priority-task-manager/shared/pkg/repositories"
)

var (
	tasksInQueueCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_in_queue_count",
	})

	tasksInProgressCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_in_progress_count",
	})
)

type TaskCountService struct {
	inQueueTaskCountRepository    repositories.NoKeyReader[uint]
	inProgressTaskCountRepository repositories.NoKeyReader[uint]
}

func MakeTaskCountService(
	inQueueTaskCountRepository repositories.NoKeyReader[uint],
	inProgressTaskCountRepository repositories.NoKeyReader[uint],
) TaskCountService {
	return TaskCountService{
		inQueueTaskCountRepository:    inQueueTaskCountRepository,
		inProgressTaskCountRepository: inProgressTaskCountRepository,
	}
}

func (tcs TaskCountService) UpdateQueued() error {
	currentTasksInQueueCount, err := tcs.inQueueTaskCountRepository.Get()
	if err != nil {
		return err
	}

	tasksInQueueCount.Set(float64(currentTasksInQueueCount))

	return nil
}

func (tcs TaskCountService) UpdateInProgress() error {
	currentTasksInProgressCount, err := tcs.inProgressTaskCountRepository.Get()
	if err != nil {
		return err
	}

	tasksInProgressCount.Set(float64(currentTasksInProgressCount))

	return nil
}
