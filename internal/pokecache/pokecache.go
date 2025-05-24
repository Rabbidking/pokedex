package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu    sync.Mutex
	cache map[string]cacheEntry
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		mu:    sync.Mutex{},
		cache: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, value []byte) {
	//Add entry to cache

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       value,
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	//Get an entry from the cache

	c.mu.Lock()
	defer c.mu.Unlock()

	value, found := c.cache[key]
	return value.val, found

}

func (c *Cache) reapLoop(interval time.Duration) {
	//removes cache entries that are older than the time interval passed into it from NewCache
	ticker := time.NewTicker(interval)

	c.mu.Lock()
	defer c.mu.Unlock()

	for range ticker.C {

		for key, val := range c.cache {
			if time.Now().Sub(val.createdAt) > interval {
				delete(c.cache, key)
			}
		}
	}
}
