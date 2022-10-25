package repositories

import "priority-task-manager/shared/pkg/types"

type (
	Reader[Model any, Key comparable] interface {
		Get(key Key) (Model, error)
	}

	Repository[Model any, Key comparable] interface {
		Reader[Model, Key]
		Add(entity Model) error
		Update(entity Model) error
		Delete(id types.ID) error
	}
)
