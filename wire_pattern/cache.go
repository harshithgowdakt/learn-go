package main

import "time"

// Cache
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
}

type RedisCache struct {
	// redis client would go here
}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
	// Redis implementation
	return nil, false
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) {
	// Redis implementation
}
