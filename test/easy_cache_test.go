package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/LittleMoreInteresting/cache"
	"github.com/LittleMoreInteresting/cache/fast"
	"github.com/LittleMoreInteresting/cache/lru"
)

/*func TestEasyCacheGet(t *testing.T) {
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

}*/
func BenchmarkNewCacheSetParallel(b *testing.B) {
	newEasyCache := cache.NewEasyCache(lru.New(b.N*100, nil), nil)
	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			newEasyCache.Set(parallelKey(id, counter), string(value()))
			counter = counter + 1
		}
	})
}
func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}
func value() []byte {
	return make([]byte, 100)
}
func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

func BenchmarkTourFastCacheSetParallel(b *testing.B) {
	newFastCache := fast.NewFastCache(b.N, 64, nil)
	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			newFastCache.Set(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}

//go test -bench=SetParallel -benchmem -count=10 ./... | tee old.txt
// benchstat old.txt
