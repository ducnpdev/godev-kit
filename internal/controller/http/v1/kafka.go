package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProducerRequest godoc
// @Summary      Send a message to a Kafka topic
// @Description  Send a message to a Kafka topic
// @Tags         kafka
// @Accept       json
// @Produce      json
// @Param        request body struct{topic string; key string; value interface{}} true "Kafka message"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/kafka/producer/request [post]
func (h *V1) ProducerRequest(c *gin.Context) {
	var req struct {
		Topic string      `json:"topic" binding:"required"`
		Key   string      `json:"key"`
		Value interface{} `json:"value" binding:"required"`
	}
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
