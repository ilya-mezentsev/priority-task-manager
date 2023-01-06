package stat

import "github.com/jmoiron/sqlx"

//goland:noinspection all
const (
	statExistenceQuery = `select count(*) > 0 from task_stat`
)

type ExistenceRepository struct {
	db *sqlx.DB
}

func MakeExistenceRepository(db *sqlx.DB) ExistenceRepository {
	return ExistenceRepository{db: db}
}

func (e ExistenceRepository) Get() (bool, error) {
	var statExists bool
	err := e.db.Get(&statExists, statExistenceQuery)

	return statExists, err
}
