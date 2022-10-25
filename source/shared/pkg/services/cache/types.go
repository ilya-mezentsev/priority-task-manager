package cache

type Cached[Model any] struct {
	value    Model
	cachedTs int64
}
