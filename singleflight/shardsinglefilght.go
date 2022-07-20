package singleflight

import "hash/fnv"

type shardSingleFlight struct {
	shards []SingleFlight
}

const SHARDS = 64

func NewShardSingleFlight() *shardSingleFlight {
	ssf := &shardSingleFlight{
		shards: make([]SingleFlight, SHARDS),
	}
	for i := 0; i < SHARDS; i++ {
		ssf.shards[i] = NewSingleFlight()
	}
	return ssf
}

func (ssf *shardSingleFlight) getShard(key string) SingleFlight {
	hashKey := fnv32(key)
	return ssf.shards[hashKey&(SHARDS-1)]
}

func (ssf *shardSingleFlight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	sf := ssf.getShard(key)
	return sf.Do(key, fn)
}

func (ssf *shardSingleFlight) DoEx(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	sf := ssf.getShard(key)
	return sf.DoEx(key, fn)
}

func fnv32(key string) uint32 {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(key))
	hash.Sum32()
	return hash.Sum32()
}
