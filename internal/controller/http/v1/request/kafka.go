package request

// KafkaMessage represents a message to be sent to Kafka.
type KafkaMessage struct {
	Topic string      `json:"topic" binding:"required"`
	Key   string      `json:"key"`
	Value interface{} `json:"value" binding:"required"`
}
