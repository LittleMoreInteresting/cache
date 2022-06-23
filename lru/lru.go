package lru

import (
	"container/list"

	"github.com/LittleMoreInteresting/cache"
)

// lru Cache Not concurrent security
type lru struct {
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
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - cache.CalcLen(en.value) +
			cache.CalcLen(value)
		en.value = value
		return
	}
	en := &entry{key: key, value: value}
	e := l.ll.PushBack(en)
	l.cache[key] = e
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	l.ll.Remove(e)
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)
	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

func (l *lru) Len() int {
	return l.ll.Len()
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
