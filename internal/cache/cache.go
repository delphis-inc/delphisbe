package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

const KEY_PREFIX = "v1"

type ChathamCache interface {
	Set(key string, val interface{}, dur time.Duration)
	Get(key string) (interface{}, bool)
}

type inMemoryCache struct {
	c         *goCache.Cache
	keyPrefix string
}

func NewInMemoryCache() ChathamCache {
	return &inMemoryCache{
		c:         goCache.New(time.Hour, 2*time.Hour),
		keyPrefix: KEY_PREFIX,
	}
}

func (c *inMemoryCache) Set(key string, val interface{}, dur time.Duration) {
	c.c.Set(key, val, dur)
}

func (c *inMemoryCache) Get(key string) (interface{}, bool) {
	return c.c.Get(key)
}
