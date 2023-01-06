package services

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	repositoriesConstructor "priority-task-manager/metrics/internal/repositories"
	"priority-task-manager/metrics/internal/services/metrics"
	"priority-task-manager/shared/pkg/types"
	"sync"
	"time"
)

type Services struct {
	uniqueRoles             []types.Role
	taskCountMetricsService metrics.TaskCountService
	queueWaitingTimeService metrics.WaitingTimeService
}

func MakeServices(db *sqlx.DB) Services {
	repositories := repositoriesConstructor.MakeRepositories(db)

	uniqueRolesRepository := repositories.UniqueRolesRepository()
	uniqueRoles, err := uniqueRolesRepository.Get()
	if err != nil {
		log.Fatalf("Unable to get unique roles: %v", err)
	}

	return Services{
		uniqueRoles: uniqueRoles,
		taskCountMetricsService: metrics.MakeTaskCountService(
			repositories.GeneralTaskCountRepository(),
			repositories.InQueueTaskCountRepository(),
			repositories.InProgressTaskCountRepository(),
			repositories.CompletedCountRepository(),
		),

		queueWaitingTimeService: metrics.MakeWaitingTimeService(
			repositories.AvgExtractedWaitingTimeRepository(),
			repositories.AvgQueuedWaitingTimeRepository(),
			repositories.AvgCompletedWaitingTimeRepository(),
		),
	}
}

func (ss Services) StartObserveMetrics() {
	targets := map[string]func(role types.Role) error{
		"update_general_tasks_count":       ss.taskCountMetricsService.UpdateGeneral,
		"update_in_progress_tasks_count":   ss.taskCountMetricsService.UpdateInProgress,
		"update_queued_tasks_count":        ss.taskCountMetricsService.UpdateQueued,
		"update_completed_tasks_count":     ss.taskCountMetricsService.UpdateCompleted,
		"update_extracting_waiting_time":   ss.queueWaitingTimeService.UpdateExtracted,
		"update_in_queue_waiting_time":     ss.queueWaitingTimeService.UpdateQueued,
		"update_in_completed_waiting_time": ss.queueWaitingTimeService.UpdateComplete,
	}

	for {
		var wg sync.WaitGroup
		for key, target := range targets {
			for _, role := range ss.uniqueRoles {
				wg.Add(1)

				go func(key string, role types.Role, target func(role types.Role) error) {
					defer wg.Done()

					err := target(role)
					if err != nil {
						log.Errorf("Unable to exec %s: %v\n", key, err)
					}
				}(key, role, target)
			}
		}

		wg.Wait()

		time.Sleep(time.Second * 2)
	}
}
