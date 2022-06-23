package fifo

import (
	"container/list"

	"github.com/LittleMoreInteresting/cache"
)

// FIFO Cache Not concurrent security
type fifo struct {
	// Maximum cache capacity .bytes
	maxBytes int
	//This callback function is called when an entry is removed from the cache.
	//The default value is ni
	onEvicted func(key string, value interface{})
	//The number of bytes used, including only the value, not the key
	usedBytes int

	ll *list.List

	cache map[string]*list.Element
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value) +
			cache.CalcLen(value)
		en.value = value
		return
	}
	en := &entry{key: key, value: value}
	e := f.ll.PushBack(en)
	f.cache[key] = e
	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}
	return nil
}

func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
