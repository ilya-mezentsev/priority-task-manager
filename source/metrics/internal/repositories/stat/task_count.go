package stat

import "github.com/jmoiron/sqlx"

const (
	inQueueCountQuery    = `select count(*) from task_stat where extracted_from_queue is null`
	inProgressCountQuery = `select count(*) from task_stat where extracted_from_queue is not null and completed is null`
	completedCountQuery  = `select count(*) from task_stat where completed is not null`
)

type TaskCountRepository struct {
	db    *sqlx.DB
	query string
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

func (tc TaskCountRepository) Get() (uint, error) {
	var count uint
	err := tc.db.Get(&count, tc.query)

	return count, err
}
