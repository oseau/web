// Package redis contains the redis connection and the models for the redis.
package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
)

// Redis is the redis connection
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new redis connection
func NewRedis() *Redis {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "redis:6379", // we use redis in docker compose
			Password: "",
			DB:       0,
		})
	})
	return &Redis{client: client}
}

// Close closes the redis connection
func (r *Redis) Close() error {
	return r.client.Close()
}
