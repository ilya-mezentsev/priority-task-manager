package stat

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

//goland:noinspection all
const (
	avgExtractedWaitingTimeQuery = `
		select
			avg(
				extract(epoch from (extracted_from_queue::timestamp - added_to_queue::timestamp))
			)
		from task_stat
		where
			extracted_from_queue is not null and
			(select role from account where hash = account_hash) = $1
	`

	avgQueuedWaitingTimeQuery = `
		select
			avg(
				extract(epoch from (current_timestamp - added_to_queue::timestamp))
			)
		from task_stat
		where (select role from account where hash = account_hash) = $1
	`

	avgCompleteWaitingTimeQuery = `
		select
			avg(
				extract(epoch from (completed::timestamp - extracted_from_queue::timestamp))
			)
		from task_stat
		where
			extracted_from_queue is not null and
			completed is not null and
			(select role from account where hash = account_hash) = $1
	`
)

type AvgWaitingTimeRepository struct {
	db    *sqlx.DB
	query string
}

func MakeAvgExtractedWaitingTimeRepository(db *sqlx.DB) AvgWaitingTimeRepository {
	return AvgWaitingTimeRepository{
		db:    db,
		query: avgExtractedWaitingTimeQuery,
	}
}

func MakeAvgQueuedWaitingTimeRepository(db *sqlx.DB) AvgWaitingTimeRepository {
	return AvgWaitingTimeRepository{
		db:    db,
		query: avgExtractedWaitingTimeQuery,
	}
}

func MakeAvgCompleteWaitingTimeRepository(db *sqlx.DB) AvgWaitingTimeRepository {
	return AvgWaitingTimeRepository{
		db:    db,
		query: avgCompleteWaitingTimeQuery,
	}
}

func (a AvgWaitingTimeRepository) Get(key types.Role) (float64, error) {
	var avgWaitingTime float64
	err := a.db.Get(&avgWaitingTime, a.query, key)

	return avgWaitingTime, err
}
