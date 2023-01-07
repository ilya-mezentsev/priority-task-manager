package services

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	repositoriesConstructor "priority-task-manager/metrics/internal/repositories"
	"priority-task-manager/metrics/internal/services/metrics"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"sync"
	"time"
)

type Services struct {
	uniqueRoles             []types.Role
	statExistenceRepository repositories.NoKeyReader[bool]
	taskCountMetricsService metrics.TaskCountService
	queueWaitingTimeService metrics.WaitingTimeService
}

func MakeServices(db *sqlx.DB) Services {
	repos := repositoriesConstructor.MakeRepositories(db)

	uniqueRolesRepository := repos.UniqueRolesRepository()
	uniqueRoles, err := uniqueRolesRepository.Get()
	if err != nil {
		log.Fatalf("Unable to get unique roles: %v", err)
	}

	return Services{
		uniqueRoles:             uniqueRoles,
		statExistenceRepository: repos.StatExistenceRepository(),
		taskCountMetricsService: metrics.MakeTaskCountService(
			repos.GeneralTaskCountRepository(),
			repos.InQueueTaskCountRepository(),
			repos.InProgressTaskCountRepository(),
			repos.CompletedCountRepository(),
		),

		queueWaitingTimeService: metrics.MakeWaitingTimeService(
			repos.AvgExtractedWaitingTimeRepository(),
			repos.AvgQueuedWaitingTimeRepository(),
			repos.AvgCompletedWaitingTimeRepository(),
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
		statExists, err := ss.statExistenceRepository.Get()
		if err != nil {
			log.Errorf("Unable to check stat existence: %v", err)

			wait()
			continue
		}

		if !statExists {
			log.Info("No stats for tasks, skipping metrics observing")

			wait()
			continue
		}

		var wg sync.WaitGroup
		for key, target := range targets {
			for _, role := range ss.uniqueRoles {
				wg.Add(1)

				go func(key string, role types.Role, target func(role types.Role) error) {
					defer wg.Done()

					metricsErr := target(role)
					if metricsErr != nil {
						log.Errorf("Unable to exec %s: %v\n", key, metricsErr)
					}
				}(key, role, target)
			}
		}

		wg.Wait()

		wait()
	}
}

func wait() {
	// просто не хочется копиастить эти 2 секунды в 3 места
	time.Sleep(time.Second * 2)
}
