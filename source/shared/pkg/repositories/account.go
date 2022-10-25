package repositories

import (
	"github.com/jmoiron/sqlx"
	"priority-task-manager/shared/pkg/types"
)

//goland:noinspection SqlNoDataSourceInspection
const getAccountQuery = `select hash, role from account where hash = $1`

type Account struct {
	db *sqlx.DB
}

func MakeAccountRepository(db *sqlx.DB) Account {
	return Account{
		db: db,
	}
}

func (a Account) Get(key string) (types.Account, error) {
	var model types.Account
	err := a.db.Get(&model, getAccountQuery, key)

	return model, err
}
