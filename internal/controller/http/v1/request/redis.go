package request

// RedisValue represents a key-value pair for Redis.
type RedisValue struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
