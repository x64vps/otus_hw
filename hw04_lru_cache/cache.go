package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.RWMutex
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.Lock()
	defer lru.Unlock()

	item := cacheItem{key: key, value: value}

	if l, ok := lru.items[key]; ok {
		l.Value = item
		lru.queue.MoveToFront(l)

		return true
	}

	lru.items[key] = lru.queue.PushFront(item)

	if lru.queue.Len() > lru.capacity {
		lru.removeLast()
	}

	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.RLock()
	defer lru.RUnlock()

	if _, ok := lru.items[key]; !ok {
		return nil, false
	}

	item := lru.items[key]
	lru.queue.MoveToFront(item)

	return item.Value.(cacheItem).value, true
}

func (lru *lruCache) Clear() {
	lru.Lock()
	defer lru.Unlock()

	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.queue = NewList()
}

func (lru *lruCache) removeLast() {
	last := lru.queue.Back()

	delete(lru.items, last.Value.(cacheItem).key)

	lru.queue.Remove(last)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
