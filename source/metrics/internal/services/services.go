package services

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	repositoriesConstructor "priority-task-manager/metrics/internal/repositories"
	"priority-task-manager/metrics/internal/services/metrics"
	"time"
)

type Services struct {
	taskCountMetricsService metrics.TaskCountService
}

func MakeServices(db *sqlx.DB) Services {
	repositories := repositoriesConstructor.MakeRepositories(db)

	return Services{
		taskCountMetricsService: metrics.MakeTaskCountService(
			repositories.InQueueTaskCountRepository(),
			repositories.InProgressTaskCountRepository(),
		),
	}
}

func (ss Services) StartObserveMetrics() {
	targets := map[string]func() error{
		"update_in_progress_tasks_count": ss.taskCountMetricsService.UpdateInProgress,
		"update_queued_tasks_count":      ss.taskCountMetricsService.UpdateQueued,
	}

	for {
		for key, target := range targets {
			err := target()
			if err != nil {
				log.Errorf("Unable to exec %s: %v\n", key, err)
			}
		}

		time.Sleep(time.Second * 2)
	}
}
