package lfu

import (
	"testing"

	"github.com/matryer/is"
)

func TestSet(t *testing.T) {
	Is := is.New(t)
	cache := New(64, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	Is.Equal(v, 1)

}

func TestOnEvicted(t *testing.T) {
	isObj := is.New(t)
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(64, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	cache.Get("k1")
	cache.Set("k3", 3)
	cache.Set("k4", 4)
	cache.Del("k1")
	cache.Del("k2")
	expected := []string{"k1", "k2"}
	isObj.Equal(expected, keys)
	isObj.Equal(2, cache.Len())
}
