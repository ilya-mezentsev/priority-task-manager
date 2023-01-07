package cache

import (
	"priority-task-manager/shared/pkg/repositories"
	"priority-task-manager/shared/pkg/types"
	"sync"
	"time"
)

type Service[Model any, Key comparable] struct {
	sync.Mutex

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

	cached, found := r.get(key)
	if found {
		// кеш еще не протух
		if int(cached.cachedTs-currentTs) < r.cacheLifetimeSeconds {
			return cached.value, nil
		} else {
			r.delete(key)
			return r.Get(key)
		}
	} else {
		model, err := r.source.Get(key)
		if err != nil {
			return zeroResult, err
		}

		r.set(key, Cached[Model]{
			value:    model,
			cachedTs: currentTs,
		})

		return model, nil
	}
}

func (r *Service[Model, Key]) set(key Key, value Cached[Model]) {
	r.Lock()
	defer r.Unlock()

	r.cache[key] = value
}

func (r *Service[Model, Key]) get(key Key) (Cached[Model], bool) {
	r.Lock()
	defer r.Unlock()

	value, found := r.cache[key]

	return value, found
}

func (r *Service[Model, Key]) delete(key Key) {
	r.Lock()
	defer r.Unlock()

	delete(r.cache, key)
}

func (r *Service[Model, Key]) Add(entity Model) error {
	r.clear()
	return r.repository().Add(entity)
}

func (r *Service[Model, Key]) clear() {
	r.Lock()
	defer r.Unlock()

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
