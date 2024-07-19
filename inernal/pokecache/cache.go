package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]entry
	mu sync.RWMutex
}

type entry struct {
	createdAt time.Time
	val []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = entry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cache, ok := c.cache[key]

	if !ok {
		return nil, false
	}
	return cache.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.cache {
			if diff := time.Since(val.createdAt); diff >= interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		cache: make(map[string]entry),
	}

	go newCache.reapLoop(interval)

	return newCache
}