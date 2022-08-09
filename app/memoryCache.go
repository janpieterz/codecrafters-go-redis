package main

import "time"

type CacheItem struct {
	value      string
	expiration time.Time
}

type MemoryCache map[string]CacheItem

func (self MemoryCache) Push(key string, value string, expiration time.Time) {
	self[key] = CacheItem{value: value, expiration: expiration}
}

func (self MemoryCache) Get(key string) string {
	value, valueExisted := self[key]
	if valueExisted {
		if value.expiration.Before(time.Now()) {
			delete(self, key)
		} else {
			return value.value

		}
	}
	return "UNKNOWN"
}
