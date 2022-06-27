package cache

type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

type EasyCache struct {
	mainCache *safeCache
	getter    Getter
}

func NewEasyCache(cache Cache, getter Getter) *EasyCache {
	return &EasyCache{mainCache: newSafeCache(cache), getter: getter}
}

func (cache EasyCache) Get(key string) interface{} {
	value := cache.mainCache.get(key)
	if value != nil {
		return value
	}
	if cache.getter != nil {
		value := cache.getter.Get(key)
		if value == nil {
			return nil
		}
		cache.mainCache.set(key, value)
		return value
	}
	return nil
}

func (cache EasyCache) Set(key string, value interface{}) {
	cache.mainCache.set(key, value)
}

func (cache EasyCache) Stat() *Stat {
	return cache.mainCache.stat()
}
