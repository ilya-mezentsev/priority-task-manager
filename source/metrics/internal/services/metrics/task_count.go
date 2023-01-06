package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
)

var (
	generalTasksCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "general_tasks_count",
	}, []string{"role"})

	tasksInQueueCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tasks_in_queue_count",
	}, []string{"role"})

	tasksInProgressCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tasks_in_progress_count",
	}, []string{"role"})

	completedTaskCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "completed_tasks_count",
	}, []string{"role"})
)

type TaskCountService struct {
	generalTaskCountRepository    repositories.Reader[uint, types.Role]
	inQueueTaskCountRepository    repositories.Reader[uint, types.Role]
	inProgressTaskCountRepository repositories.Reader[uint, types.Role]
	completedTaskCountRepository  repositories.Reader[uint, types.Role]
}

func MakeTaskCountService(
	generalTaskCountRepository repositories.Reader[uint, types.Role],
	inQueueTaskCountRepository repositories.Reader[uint, types.Role],
	inProgressTaskCountRepository repositories.Reader[uint, types.Role],
	completedTaskCountRepository repositories.Reader[uint, types.Role],
) TaskCountService {
	return TaskCountService{
		generalTaskCountRepository:    generalTaskCountRepository,
		inQueueTaskCountRepository:    inQueueTaskCountRepository,
		inProgressTaskCountRepository: inProgressTaskCountRepository,
		completedTaskCountRepository:  completedTaskCountRepository,
	}
}

func (tcs TaskCountService) UpdateGeneral(role types.Role) error {
	generalTaskCount, err := tcs.generalTaskCountRepository.Get(role)
	if err != nil {
		return err
	}

	generalTasksCount.WithLabelValues(string(role)).Set(float64(generalTaskCount))

	return nil
}

func (tcs TaskCountService) UpdateQueued(role types.Role) error {
	currentTasksInQueueCount, err := tcs.inQueueTaskCountRepository.Get(role)
	if err != nil {
		return err
	}

	tasksInQueueCount.WithLabelValues(string(role)).Set(float64(currentTasksInQueueCount))

	return nil
}

func (tcs TaskCountService) UpdateInProgress(role types.Role) error {
	currentTasksInProgressCount, err := tcs.inProgressTaskCountRepository.Get(role)
	if err != nil {
		return err
	}

	tasksInProgressCount.WithLabelValues(string(role)).Set(float64(currentTasksInProgressCount))

	return nil
}

func (tcs TaskCountService) UpdateCompleted(role types.Role) error {
	currentCompletedTaskCount, err := tcs.completedTaskCountRepository.Get(role)
	if err != nil {
		return err
	}

	completedTaskCount.WithLabelValues(string(role)).Set(float64(currentCompletedTaskCount))

	return nil
}
