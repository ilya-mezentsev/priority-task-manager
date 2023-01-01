package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/metrics/internal/repositories/stat"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
)

type Repositories struct {
	inQueueTaskCountRepository        repositories.NoKeyReader[uint]
	inProgressTaskCountRepository     repositories.NoKeyReader[uint]
	uniqueRolesRepository             repositories.NoKeyReader[[]types.Role]
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role]
	avgQueuedWaitingTimeRepository    repositories.Reader[float64, types.Role]
}

func MakeRepositories(db *sqlx.DB) Repositories {
	return Repositories{
		inQueueTaskCountRepository:        stat.MakeInQueueCountRepository(db),
		inProgressTaskCountRepository:     stat.MakeInProgressCountRepository(db),
		uniqueRolesRepository:             MakeUniqueRolesRepository(db),
		avgExtractedWaitingTimeRepository: stat.MakeAvgExtractedWaitingTimeRepository(db),
		avgQueuedWaitingTimeRepository:    stat.MakeAvgQueuedWaitingTimeRepository(db),
	}
}

func (r Repositories) InQueueTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inQueueTaskCountRepository
}

func (r Repositories) InProgressTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inProgressTaskCountRepository
}

func (r Repositories) UniqueRolesRepository() repositories.NoKeyReader[[]types.Role] {
	return r.uniqueRolesRepository
}

func (r Repositories) AvgExtractedWaitingTimeRepository() repositories.Reader[float64, types.Role] {
	return r.avgExtractedWaitingTimeRepository
}

func (r Repositories) AvgQueuedWaitingTimeRepository() repositories.Reader[float64, types.Role] {
	return r.avgQueuedWaitingTimeRepository
}
