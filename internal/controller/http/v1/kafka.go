package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/gin-gonic/gin"
)

// ProducerRequest godoc
// @Summary      Send a message to a Kafka topic
// @Description  Send a message to a Kafka topic
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Param        request body request.KafkaMessage true "Kafka message"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/kafka/producer/request [post]
func (h *V1) ProducerRequest(c *gin.Context) {
	var req request.KafkaMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.kafka.ProduceMessage(c.Copy(), req.Topic, req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

// ConsumerReceiver godoc
// @Summary      Receive a message from a Kafka topic and group
// @Description  Receive a message from a Kafka topic and group
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Param        topic query string true "Kafka topic"
// @Param        group query string true "Kafka group"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      504  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/kafka/consumer/receiver [get]
func (h *V1) ConsumerReceiver(c *gin.Context) {
	topic := c.Query("topic")
	group := c.Query("group")
	if topic == "" || group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic and group are required"})
		return
	}

	msgCh := make(chan map[string]interface{}, 1)
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	a, b, err := h.kafka.ConsumeMessage(ctx, topic, group)
	// if err := nil {
	fmt.Println(a, b)
	// }
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start consumer (will exit after one message)
	// go h.KafkaRepo.StartConsumer(ctx, topic)

	select {
	case msg := <-msgCh:
		c.JSON(http.StatusOK, msg)
	case <-ctx.Done():
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "no message received"})
	}
}

// EnableProducer godoc
// @Summary      Enable Kafka producer
// @Description  Enable the Kafka producer to send messages
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /v1/kafka/producer/enable [post]
func (h *V1) EnableProducer(c *gin.Context) {
	h.kafka.EnableProducer()
	c.JSON(http.StatusOK, gin.H{"status": "producer enabled"})
}

// DisableProducer godoc
// @Summary      Disable Kafka producer
// @Description  Disable the Kafka producer from sending messages
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /v1/kafka/producer/disable [post]
func (h *V1) DisableProducer(c *gin.Context) {
	h.kafka.DisableProducer()
	c.JSON(http.StatusOK, gin.H{"status": "producer disabled"})
}

// EnableConsumer godoc
// @Summary      Enable Kafka consumer
// @Description  Enable the Kafka consumer to receive messages
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /v1/kafka/consumer/enable [post]
func (h *V1) EnableConsumer(c *gin.Context) {
	h.kafka.EnableConsumer()
	c.JSON(http.StatusOK, gin.H{"status": "consumer enabled"})
}

// DisableConsumer godoc
// @Summary      Disable Kafka consumer
// @Description  Disable the Kafka consumer from receiving messages
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /v1/kafka/consumer/disable [post]
func (h *V1) DisableConsumer(c *gin.Context) {
	h.kafka.DisableConsumer()
	c.JSON(http.StatusOK, gin.H{"status": "consumer disabled"})
}

// GetKafkaStatus godoc
// @Summary      Get Kafka status
// @Description  Get the current status of Kafka producer and consumer
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/kafka/status [get]
func (h *V1) GetKafkaStatus(c *gin.Context) {
	status := h.kafka.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   status,
	})
}

// CheckConsumerStatus godoc
// @Summary      Check consumer status
// @Description  Check if Kafka consumer is enabled or disabled
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/kafka/consumer/status [get]
func (h *V1) CheckConsumerStatus(c *gin.Context) {
	isEnabled := h.kafka.IsConsumerEnabled()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"consumer_enabled": isEnabled,
			"message":          getConsumerStatusMessage(isEnabled),
		},
	})
}

// CheckProducerStatus godoc
// @Summary      Check producer status
// @Description  Check if Kafka producer is enabled or disabled
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /v1/kafka/producer/status [get]
func (h *V1) CheckProducerStatus(c *gin.Context) {
	isEnabled := h.kafka.IsProducerEnabled()
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"producer_enabled": isEnabled,
			"message":          getProducerStatusMessage(isEnabled),
		},
	})
}

// getConsumerStatusMessage returns a human-readable message for consumer status
func getConsumerStatusMessage(enabled bool) string {
	if enabled {
		return "Kafka consumer is enabled and can receive messages"
	}
	return "Kafka consumer is disabled and cannot receive messages"
}

// getProducerStatusMessage returns a human-readable message for producer status
func getProducerStatusMessage(enabled bool) string {
	if enabled {
		return "Kafka producer is enabled and can send messages"
	}
	return "Kafka producer is disabled and cannot send messages"
}
