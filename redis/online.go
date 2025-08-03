package redis

import (
	"context"
)

const (
	keyViewCount = "view-count"
)

// AddViewCount adds a new view count, and returns the number of views
func (r *Redis) AddViewCount() int64 {
	return r.client.Incr(context.Background(), keyViewCount).Val()
}
