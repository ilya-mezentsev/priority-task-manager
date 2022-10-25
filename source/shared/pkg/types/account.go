package types

type Account struct {
	Hash string `db:"hash"`
	Role string `db:"role"`
}
