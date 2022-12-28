package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/metrics/internal/repositories/stat"
	"priority-task-manager/shared/pkg/repositories"
)

type Repositories struct {
	inQueueTaskCountRepository    repositories.NoKeyReader[uint]
	inProgressTaskCountRepository repositories.NoKeyReader[uint]
}

func MakeRepositories(db *sqlx.DB) Repositories {
	return Repositories{
		inQueueTaskCountRepository:    stat.MakeInQueueCountRepository(db),
		inProgressTaskCountRepository: stat.MakeInProgressCountRepository(db),
	}
}

func (r Repositories) InQueueTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inQueueTaskCountRepository
}

func (r Repositories) InProgressTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inProgressTaskCountRepository
}
