package db

import (
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	cc *cache.Cache
}

func NewCache() *Cache {
	return &Cache{cache.New(time.Hour, time.Hour*24)}
}

func (c *Cache) Set(k string, v string) {
	c.cc.Set(k, v, cache.DefaultExpiration)
}

func (c *Cache) Get(k string) (string, bool) {
	v, ok := c.cc.Get(k)
	if !ok {
		return "", false
	}
	str, ok := v.(string)
	if !ok {
		log.Printf("Found a non-string in cache `%#v`", v)
		return "", false
	}
	return str, ok
}
