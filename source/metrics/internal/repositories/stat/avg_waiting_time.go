package stat

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

//goland:noinspection all
const avgWaitingTimeQuery = `
	select
		avg(
			extract(epoch from (extracted_from_queue::timestamp - added_to_queue::timestamp))
		)
	from task_stat
	where
		extracted_from_queue is not null and
		(select role from account where hash = account_hash) = $1
`

type AvgWaitingTimeRepository struct {
	db *sqlx.DB
}

func MakeAvgWaitingTimeRepository(db *sqlx.DB) AvgWaitingTimeRepository {
	return AvgWaitingTimeRepository{db}
}

func (a AvgWaitingTimeRepository) Get(key types.Role) (float64, error) {
	var avgWaitingTime float64
	err := a.db.Get(&avgWaitingTime, avgWaitingTimeQuery, key)

	return avgWaitingTime, err
}
