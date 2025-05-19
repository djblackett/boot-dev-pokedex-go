package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	cache map[string]CacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) *Cache {

	cacheMap := make(map[string]CacheEntry)
	cache := &Cache{
		cache: cacheMap,
		mu:    sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	currentTime := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CacheEntry{
		CreatedAt: currentTime,
		Val:       val,
	}
}

func (c *Cache) Get(key string) (CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.cache[key]
	if !ok {
		return CacheEntry{}, false
	}
	return i, true
}

func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.mu.Lock()
		for key, val := range c.cache {
			if val.CreatedAt.Before(t.Add(-interval)) {

				delete(c.cache, key)

			}
		}
		c.mu.Unlock()
	}
}
