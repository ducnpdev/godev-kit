package entity

// RedisValue represents a key-value pair for Redis.
type RedisValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
