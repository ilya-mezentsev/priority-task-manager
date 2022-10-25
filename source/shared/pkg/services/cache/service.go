package cache

import (
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"time"
)

type Service[Model any, Key comparable] struct {
	cacheLifetimeSeconds int
	cache                map[Key]Cached[Model]

	source repositories.Reader[Model, Key]
}

func MakeService[Model any, Key comparable](source repositories.Reader[Model, Key], cacheLifetimeSeconds int) *Service[Model, Key] {
	return &Service[Model, Key]{
		cacheLifetimeSeconds: cacheLifetimeSeconds,
		cache:                map[Key]Cached[Model]{},
		source:               source,
	}
}

func (r *Service[Model, Key]) Get(key Key) (Model, error) {
	var zeroResult Model

	now := time.Now()
	currentTs := now.Unix()

	cached, found := r.cache[key]
	if found {
		// кеш еще не протух
		if int(cached.cachedTs-currentTs) < r.cacheLifetimeSeconds {
			return cached.value, nil
		} else {
			delete(r.cache, key)
			return r.Get(key)
		}
	} else {
		model, err := r.source.Get(key)
		if err != nil {
			return zeroResult, err
		}

		r.cache[key] = Cached[Model]{
			value:    model,
			cachedTs: currentTs,
		}

		return model, nil
	}
}

func (r *Service[Model, Key]) Add(entity Model) error {
	r.clear()
	return r.repository().Add(entity)
}

func (r *Service[Model, Key]) clear() {
	r.cache = map[Key]Cached[Model]{}
}

func (r *Service[Model, Key]) repository() repositories.Repository[Model, Key] {
	// Делаем каст тут, чтобы не мучаться.
	// Сейчас сложно представить ситуацию, в которой этот код упадет.
	return r.source.(repositories.Repository[Model, Key])
}

func (r *Service[Model, Key]) Update(entity Model) error {
	r.clear()
	return r.repository().Update(entity)
}

func (r *Service[Model, Key]) Delete(id types.ID) error {
	r.clear()
	return r.repository().Delete(id)
}
