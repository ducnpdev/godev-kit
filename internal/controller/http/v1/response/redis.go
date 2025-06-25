package response

import "github.com/ducnpdev/godev-kit/internal/entity"

// RedisValue represents a key-value pair for Redis.
type RedisValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewRedisValue creates a new RedisValue response from an entity.
func NewRedisValue(v entity.RedisValue) RedisValue {
	return RedisValue{
		Key:   v.Key,
		Value: v.Value,
	}
}
