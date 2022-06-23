package lfu

import (
	"container/heap"

	"github.com/LittleMoreInteresting/cache"
)

type lfu struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	queue *queue
	cache map[string]*entry
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (f *lfu) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.usedBytes = f.usedBytes - cache.CalcLen(e.value) +
			cache.CalcLen(value)
		f.queue.update(e, value, e.weight+1)
		return
	}
	en := &entry{
		key:   key,
		value: value,
	}
	heap.Push(f.queue, en)
	f.cache[key] = en
	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.removeElement(heap.Pop(f.queue))
	}
}

func (f *lfu) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		f.queue.update(e, e.value, e.weight+1)
		return e.value
	}
	return nil
}

func (f *lfu) Del(key string) {
	if e, ok := f.cache[key]; ok {
		heap.Remove(f.queue, e.index)
		f.removeElement(e)
	}
}

func (f *lfu) removeElement(e interface{}) {
	if e == nil {
		return
	}
	en := e.(*entry)
	delete(f.cache, en.key)
	f.usedBytes -= en.Len()
	if f.onEvicted != nil {
		f.onEvicted(en.key, en.value)
	}
}

func (f *lfu) DelOldest() {
	if f.queue.Len() == 0 {
		return
	}
	f.removeElement(heap.Pop(f.queue))
}

func (f *lfu) Len() int {
	return f.queue.Len()
}

type entry struct {
	key    string
	value  interface{}
	weight int
	index  int
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value) + 4 + 4
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}
func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}
func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}
func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}
func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil
	en.index = -1
	*q = old[0 : n-1]
	return en
}

func (q *queue) update(en *entry, value interface{}, weight int) {
	en.value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
