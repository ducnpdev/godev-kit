package entity

import "time"

// RedisValue represents a key-value pair for Redis.
type RedisValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ShipperLocation represents a shipper's location (for Redis and DB)
type ShipperLocation struct {
	ShipperID string    `json:"shipper_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}
