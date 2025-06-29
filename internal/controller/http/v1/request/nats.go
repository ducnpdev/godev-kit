package request

// NatsPublishRequest represents the request body for NATS publish.
type NatsPublishRequest struct {
	Data string `json:"data"`
}
