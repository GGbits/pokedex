package internal

import (
	"sync"
	"time"
)

type Cache struct {
	CacheMap map[string]cacheEntry `default:"make(map[string]cacheEntry)"`
	Timeout  time.Duration         `default:"30_000_000_000"`
	mut      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(duration time.Duration) *Cache {
	c := &Cache{
		CacheMap: make(map[string]cacheEntry),
		Timeout:  duration,
	}
	ticker := time.NewTicker(duration)
	go c.reapLoop(ticker)
	return c
}

func (c *Cache) Add(str string, data []byte) {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.CacheMap[str] = cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}
}

func (c *Cache) Get(str string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()

	v, ok := c.CacheMap[str]
	if !ok {
		return nil, false
	}
	return v.val, true
}

func (c *Cache) reapLoop(tick *time.Ticker) {
	for {
		<-tick.C
		c.mut.Lock()
		for k, ce := range c.CacheMap {
			if time.Since(ce.createdAt) > c.Timeout {
				delete(c.CacheMap, k)
			}
		}
		c.mut.Unlock()
	}

}
