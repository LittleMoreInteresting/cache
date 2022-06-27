package fast

type fastCache struct {
	shards    []*cacheShard
	shardMask uint64
	hash      fnv64a
}

func NewFastCache(maxEntries int, shardsNum int, onEvicted func(key string, value interface{})) *fastCache {
	fastCache := &fastCache{
		shards:    make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
		hash:      newDefaultHasher(),
	}
	for i := 0; i < shardsNum; i++ {
		fastCache.shards[i] = newCacheShard(maxEntries, onEvicted)
	}
	return fastCache
}

func (f *fastCache) getShard(key string) *cacheShard {
	hashKey := f.hash.Sun64(key)
	return f.shards[hashKey&f.shardMask]
}

func (f *fastCache) Set(key string, value interface{}) {
	f.getShard(key).Set(key, value)
}

func (f *fastCache) Get(key string) interface{} {
	return f.getShard(key).get(key)
}

func (f *fastCache) Del(key string) {
	f.getShard(key).Del(key)
}

func (f *fastCache) Len() int {
	l := 0
	for _, s := range f.shards {
		l += s.Len()
	}
	return l
}
