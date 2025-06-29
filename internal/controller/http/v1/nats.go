package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary     Publish message
// @Description Publish a message to a NATS subject
// @ID          nats-publish
// @Tags   	    nats
// @Accept      json
// @Produce     json
// @Param       subject path string true "NATS subject"
// @Param       request body request.NatsPublishRequest true "Message data"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /v1/nats/publish/{subject} [post]
func (v1 *V1) NatsPublish(c *gin.Context) {
	subject := c.Param("subject")
	var req struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := v1.nats.Publish(subject, []byte(req.Data)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "message published"})
}

// @Summary     Subscribe to subject
// @Description Subscribe to a NATS subject (demo: returns first message)
// @ID          nats-subscribe
// @Tags   	    nats
// @Accept      json
// @Produce     json
// @Param       subject path string true "NATS subject"
// @Success     200 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /v1/nats/subscribe/{subject} [get]
func (v1 *V1) NatsSubscribe(c *gin.Context) {
	subject := c.Param("subject")
	msgCh := make(chan string, 1)
	unsubscribe, err := v1.nats.Subscribe(subject, func(msg []byte) {
		msgCh <- string(msg)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer unsubscribe()

	select {
	case msg := <-msgCh:
		c.JSON(http.StatusOK, gin.H{"message": msg})
	case <-c.Request.Context().Done():
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "no message received"})
	}
}
