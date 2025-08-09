package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	capacity int
	list     *list.List
	items    map[string]*list.Element
	mu       sync.RWMutex
}

type cacheItem struct {
	key   string
	value any
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		list:     list.New(),
		items:    make(map[string]*list.Element),
	}
}

// Getting an element and modifying the doubly linked list
func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if elem, exists := c.items[key]; exists {
		c.list.MoveToFront(elem)
		return elem.Value.(*cacheItem).value, true
	}
	return nil, false
}

// Setting a new element to the cache
func (c *LRUCache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Checking the existing of element
	if elem, exists := c.items[key]; exists {
		c.list.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
		return
	}

	// Checking the size and capacity
	if c.list.Len() == c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			c.removeElement(oldest)
		}
	}

	// Adding a new element
	item := &cacheItem{key: key, value: value}
	elem := c.list.PushFront(item)
	c.items[key] = elem
}

func (c *LRUCache) removeElement(elem *list.Element) {
	key := elem.Value.(*cacheItem).key
	delete(c.items, key)
	c.list.Remove(elem)
}
