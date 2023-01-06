package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/metrics/internal/repositories/stat"
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
)

type Repositories struct {
	generalTaskCountRepository        repositories.Reader[uint, types.Role]
	inQueueTaskCountRepository        repositories.Reader[uint, types.Role]
	inProgressTaskCountRepository     repositories.Reader[uint, types.Role]
	completedTaskCountRepository      repositories.Reader[uint, types.Role]
	avgExtractedWaitingTimeRepository repositories.Reader[float64, types.Role]
	avgQueuedWaitingTimeRepository    repositories.Reader[float64, types.Role]
	avgCompletedWaitingTimeRepository repositories.Reader[float64, types.Role]
	uniqueRolesRepository             repositories.NoKeyReader[[]types.Role]
}

func MakeRepositories(db *sqlx.DB) Repositories {
	return Repositories{
		generalTaskCountRepository:        stat.MakeGeneralCountRepository(db),
		inQueueTaskCountRepository:        stat.MakeInQueueCountRepository(db),
		inProgressTaskCountRepository:     stat.MakeInProgressCountRepository(db),
		completedTaskCountRepository:      stat.MakeCompletedCountRepository(db),
		avgExtractedWaitingTimeRepository: stat.MakeAvgExtractedWaitingTimeRepository(db),
		avgQueuedWaitingTimeRepository:    stat.MakeAvgQueuedWaitingTimeRepository(db),
		avgCompletedWaitingTimeRepository: stat.MakeAvgCompleteWaitingTimeRepository(db),
		uniqueRolesRepository:             MakeUniqueRolesRepository(db),
	}
}

func (r Repositories) GeneralTaskCountRepository() repositories.Reader[uint, types.Role] {
	return r.generalTaskCountRepository
}

func (r Repositories) InQueueTaskCountRepository() repositories.Reader[uint, types.Role] {
	return r.inQueueTaskCountRepository
}

func (r Repositories) InProgressTaskCountRepository() repositories.Reader[uint, types.Role] {
	return r.inProgressTaskCountRepository
}

func (r Repositories) CompletedCountRepository() repositories.Reader[uint, types.Role] {
	return r.completedTaskCountRepository
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

func (r Repositories) UniqueRolesRepository() repositories.NoKeyReader[[]types.Role] {
	return r.uniqueRolesRepository
}
