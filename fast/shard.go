package fast

import (
	"container/list"
	"sync"

	"github.com/LittleMoreInteresting/cache"
)

type cacheShard struct {
	lock       sync.RWMutex
	maxEntries int
	onEvicted  func(key string, value interface{})

	ll    *list.List
	cache map[string]*list.Element
}

func newCacheShard(maxEntries int, onEvicted func(key string, value interface{})) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *cacheShard) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*entry)
		en.value = value
		return
	}
	en := &entry{key: key, value: value}
	e := c.ll.PushBack(en)
	c.cache[key] = e
	if c.maxEntries > 0 && c.Len() > c.maxEntries {
		c.DelOldest()
	}
}

func (c *cacheShard) get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, ok := c.cache[key]; ok {
		c.ll.MoveToBack(v)
		return v.Value.(entry).value
	}
	return nil
}

func (c *cacheShard) Del(key string) {
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (l *cacheShard) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.cache, en.key)
	if c.onEvicted != nil {
		c.onEvicted(en.key, en.value)
	}
}

func (l *cacheShard) Len() int {
	return l.ll.Len()
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}
