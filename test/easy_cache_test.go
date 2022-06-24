package test

import (
	"sync"
	"testing"

	"github.com/LittleMoreInteresting/cache"
	"github.com/LittleMoreInteresting/cache/lru"
	"github.com/matryer/is"
)

func TestEasyCacheGet(t *testing.T) {
	mockDb := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}
	getter := cache.GetFunc(func(key string) interface{} {
		if v, ok := mockDb[key]; ok {
			return v
		}
		return nil
	})
	easyCache := cache.NewEasyCache(lru.New(64, nil), getter)
	Is := is.New(t)
	var wg sync.WaitGroup
	for key, val := range mockDb {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			Is.Equal(easyCache.Get(k), v)
			Is.Equal(easyCache.Get(k), v)
		}(key, val)
	}
	wg.Wait()
	Is.Equal(easyCache.Get("unknown"), nil)
	Is.Equal(easyCache.Get("unknown"), nil)
	stat := easyCache.Stat()
	Is.Equal(stat.Nget, 10)
	Is.Equal(stat.Nhit, 4)
}
