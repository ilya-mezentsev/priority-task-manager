package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"priority-task-manager/shared/pkg/repositories"
)

var (
	generalTasksCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "general_tasks_count",
	})

	tasksInQueueCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_in_queue_count",
	})

	tasksInProgressCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_in_progress_count",
	})

	completedTaskCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "completed_tasks_count",
	})
)

type TaskCountService struct {
	generalTaskCountRepository    repositories.NoKeyReader[uint]
	inQueueTaskCountRepository    repositories.NoKeyReader[uint]
	inProgressTaskCountRepository repositories.NoKeyReader[uint]
	completedTaskCountRepository  repositories.NoKeyReader[uint]
}

func MakeTaskCountService(
	generalTaskCountRepository repositories.NoKeyReader[uint],
	inQueueTaskCountRepository repositories.NoKeyReader[uint],
	inProgressTaskCountRepository repositories.NoKeyReader[uint],
	completedTaskCountRepository repositories.NoKeyReader[uint],
) TaskCountService {
	return TaskCountService{
		generalTaskCountRepository:    generalTaskCountRepository,
		inQueueTaskCountRepository:    inQueueTaskCountRepository,
		inProgressTaskCountRepository: inProgressTaskCountRepository,
		completedTaskCountRepository:  completedTaskCountRepository,
	}
}

func (tcs TaskCountService) UpdateGeneral() error {
	generalTaskCount, err := tcs.generalTaskCountRepository.Get()
	if err != nil {
		return err
	}

	generalTasksCount.Set(float64(generalTaskCount))

	return nil
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

func (tcs TaskCountService) UpdateCompleted() error {
	currentCompletedTaskCount, err := tcs.completedTaskCountRepository.Get()
	if err != nil {
		return err
	}

	completedTaskCount.Set(float64(currentCompletedTaskCount))

	return nil
}
