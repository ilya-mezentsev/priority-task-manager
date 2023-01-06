package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/metrics/internal/repositories/stat"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
)

type Repositories struct {
	generalTaskCountRepository        repositories.NoKeyReader[uint]
	inQueueTaskCountRepository        repositories.NoKeyReader[uint]
	inProgressTaskCountRepository     repositories.NoKeyReader[uint]
	completedTaskCountRepository      repositories.NoKeyReader[uint]
	uniqueRolesRepository             repositories.NoKeyReader[[]types.Role]
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role]
	avgQueuedWaitingTimeRepository    repositories.Reader[float64, types.Role]
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role]
}

func MakeRepositories(db *sqlx.DB) Repositories {
	return Repositories{
		generalTaskCountRepository:        stat.MakeGeneralCountRepository(db),
		inQueueTaskCountRepository:        stat.MakeInQueueCountRepository(db),
		inProgressTaskCountRepository:     stat.MakeInProgressCountRepository(db),
		completedTaskCountRepository:      stat.MakeCompletedCountRepository(db),
		uniqueRolesRepository:             MakeUniqueRolesRepository(db),
		avgExtractedWaitingTimeRepository: stat.MakeAvgExtractedWaitingTimeRepository(db),
		avgQueuedWaitingTimeRepository:    stat.MakeAvgQueuedWaitingTimeRepository(db),
		avgCompletedWaitingTimeRepository: stat.MakeAvgCompleteWaitingTimeRepository(db),
	}
}

func (r Repositories) GeneralTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.generalTaskCountRepository
}

func (r Repositories) InQueueTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inQueueTaskCountRepository
}

func (r Repositories) InProgressTaskCountRepository() repositories.NoKeyReader[uint] {
	return r.inProgressTaskCountRepository
}

func (r Repositories) CompletedCountRepository() repositories.NoKeyReader[uint] {
	return r.completedTaskCountRepository
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

func (r Repositories) AvgCompletedWaitingTimeRepository() repositories.Reader[float64, types.Role] {
	return r.avgCompletedWaitingTimeRepository
}
