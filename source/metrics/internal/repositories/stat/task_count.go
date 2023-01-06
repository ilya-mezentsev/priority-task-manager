package stat

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

const (
	generalCountQuery    = `select count(*) from task_stat where (select role from account where hash = account_hash) = $1`
	inQueueCountQuery    = `select count(*) from task_stat where extracted_from_queue is null and (select role from account where hash = account_hash) = $1`
	inProgressCountQuery = `select count(*) from task_stat where extracted_from_queue is not null and completed is null and (select role from account where hash = account_hash) = $1`
	completedCountQuery  = `select count(*) from task_stat where completed is not null and (select role from account where hash = account_hash) = $1`
)

type TaskCountRepository struct {
	db    *sqlx.DB
	query string
}

func MakeGeneralCountRepository(db *sqlx.DB) TaskCountRepository {
	return TaskCountRepository{
		db:    db,
		query: generalCountQuery,
	}
}

func MakeInQueueCountRepository(db *sqlx.DB) TaskCountRepository {
	return TaskCountRepository{
		db:    db,
		query: inQueueCountQuery,
	}
}

func MakeInProgressCountRepository(db *sqlx.DB) TaskCountRepository {
	return TaskCountRepository{
		db:    db,
		query: inProgressCountQuery,
	}
}

func MakeCompletedCountRepository(db *sqlx.DB) TaskCountRepository {
	return TaskCountRepository{
		db:    db,
		query: completedCountQuery,
	}
}

func (tc TaskCountRepository) Get(role types.Role) (uint, error) {
	var count uint
	err := tc.db.Get(&count, tc.query, string(role))

	return count, err
}
