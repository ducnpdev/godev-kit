package request

import "time"

// RedisValue represents a key-value pair for Redis.
type RedisValue struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// ShipperLocation represents a shipper's location update request
type ShipperLocation struct {
	ShipperID string    `json:"shipper_id" binding:"required"`
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"omitempty"`
}
