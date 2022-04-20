package db

import (
	"log"
	"sync"
)

// TODO: Should the cache just be a db model? because it seems like this is just not working
// TODO: Then we can just check their authId matches, and if not force a new login to reset it

// Cache basically needs the log.Printf() stuff to work so its probably a race condition I guess?
// Which is why well probably just end up with it in the db
// we can have a logout that clears cookies and clears out their record in the DB
type Cache struct {
	mutex sync.Mutex
	m     map[string]string
}

func NewCache() *Cache {
	return &Cache{mutex: sync.Mutex{}, m: make(map[string]string)}
}

func (c *Cache) Set(k string, v string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	log.Printf("k = %q, v = %q", k, v)
	c.m[k] = v
}

func (c *Cache) Get(k string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, ok := c.m[k]
	log.Printf("v = %q, ok = %t, k = %s", v, ok, k)
	if !ok {
		return "", false
	}
	return v, ok
}
