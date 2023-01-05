package services

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	repositoriesConstructor "priority-task-manager/metrics/internal/repositories"
	"priority-task-manager/metrics/internal/services/metrics"
	"sync"
	"time"
)

type Services struct {
	taskCountMetricsService metrics.TaskCountService
	queueWaitingTimeService metrics.WaitingTimeService
}

func MakeServices(db *sqlx.DB) Services {
	repositories := repositoriesConstructor.MakeRepositories(db)

	return Services{
		taskCountMetricsService: metrics.MakeTaskCountService(
			repositories.InQueueTaskCountRepository(),
			repositories.InProgressTaskCountRepository(),
			repositories.CompletedCountRepository(),
		),

		queueWaitingTimeService: metrics.MustMakeWaitingTimeService(
			repositories.UniqueRolesRepository(),
			repositories.AvgExtractedWaitingTimeRepository(),
			repositories.AvgQueuedWaitingTimeRepository(),
			repositories.AvgCompletedWaitingTimeRepository(),
		),
	}
}

func (ss Services) StartObserveMetrics() {
	targets := map[string]func() error{
		"update_in_progress_tasks_count": ss.taskCountMetricsService.UpdateInProgress,
		"update_queued_tasks_count":      ss.taskCountMetricsService.UpdateQueued,
		"update_completed_tasks_count":   ss.taskCountMetricsService.UpdateCompleted,
		"update_all_waiting_time":        ss.queueWaitingTimeService.UpdateForEachRole,
	}

	for {
		var wg sync.WaitGroup
		for key, target := range targets {
			wg.Add(1)

			go func(key string, target func() error) {
				defer wg.Done()

				err := target()
				if err != nil {
					log.Errorf("Unable to exec %s: %v\n", key, err)
				}
			}(key, target)
		}

		wg.Wait()

		time.Sleep(time.Second * 2)
	}
}
