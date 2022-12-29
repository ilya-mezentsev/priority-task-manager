package types

type (
	Role string

	Account struct {
		Hash string `db:"hash"`
		Role Role   `db:"role"`
	}
)
