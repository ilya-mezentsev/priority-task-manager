package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

//goland:noinspection all
const uniqueRolesQuery = `select distinct role from account`

type UniqueRolesRepository struct {
	db *sqlx.DB
}

func MakeUniqueRolesRepository(db *sqlx.DB) UniqueRolesRepository {
	return UniqueRolesRepository{db}
}

func (u UniqueRolesRepository) Get() ([]types.Role, error) {
	var roles []types.Role
	err := u.db.Select(&roles, uniqueRolesQuery)

	return roles, err
}
