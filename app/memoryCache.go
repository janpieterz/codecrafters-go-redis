package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	value      string
	expiration *time.Time
}

type MemoryCache struct {
	cache map[string]CacheItem
	lock  sync.RWMutex
}

func (self *MemoryCache) Push(key string, value string, expiration *time.Time) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.cache[key] = CacheItem{value: value, expiration: expiration}
}

func (self *MemoryCache) Get(key string) string {
	self.lock.Lock()
	defer self.lock.Unlock()
	value, valueExisted := self.cache[key]
	if valueExisted {
		if value.expiration != nil && value.expiration.Before(time.Now()) {
			fmt.Println("Value existed but expiration is before now")
			delete(self.cache, key)
		} else {
			return value.value

		}
	}
	return "nil"
}
