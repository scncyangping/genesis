// @Author: YangPing
// @Create: 2023/10/23
// @Description: cache配置

package pkg

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	cache *cache.Cache
}

func InitCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *Cache) Cache() *cache.Cache {
	return c.cache
}

func (c *Cache) Flush() {
	c.cache.Flush()
}
