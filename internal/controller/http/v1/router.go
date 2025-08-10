package v1

import (
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group *gin.RouterGroup, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

// NewUserRoutes -.
func NewUserRoutes(apiV1Group *gin.RouterGroup, u usecase.User, l logger.Interface) {
	r := &V1{user: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	userGroup := apiV1Group.Group("/user")

	{
		userGroup.POST("", r.CreateUser)
		userGroup.GET("", r.ListUsers)
		userGroup.GET("/:id", r.GetUser)
		userGroup.PUT("/:id", r.UpdateUser)
		userGroup.DELETE("/:id", r.DeleteUser)
	}
	apiV1Group.POST("auth/login", r.LoginUser)
}

// // NewUserRoutes -.
// func NewAuthRoutes(apiV1Group *gin.RouterGroup, u usecase.User, l logger.Interface) {
// 	r := &V1{user: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

// 	userGroup := apiV1Group.Group("/user")

// 	{
// 		userGroup.POST("", r.CreateUser)
// 		userGroup.GET("", r.ListUsers)
// 		userGroup.GET("/:id", r.GetUser)
// 		userGroup.PUT("/:id", r.UpdateUser)
// 		userGroup.DELETE("/:id", r.DeleteUser)
// 	}
// }

// NewKafkaRoutes registers Kafka producer and consumer endpoints.
func NewKafkaRoutes(apiV1Group *gin.RouterGroup, kafka usecase.Kafka, l logger.Interface) {
	r := &V1{kafka: kafka, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	kafkaGroup := apiV1Group.Group("/kafka")
	{
		kafkaGroup.POST("/producer/request", r.ProducerRequest)
		kafkaGroup.GET("/consumer/receiver", r.ConsumerReceiver)

		// Control endpoints
		kafkaGroup.POST("/producer/enable", r.EnableProducer)
		kafkaGroup.POST("/producer/disable", r.DisableProducer)
		kafkaGroup.POST("/consumer/enable", r.EnableConsumer)
		kafkaGroup.POST("/consumer/disable", r.DisableConsumer)
		kafkaGroup.GET("/status", r.GetKafkaStatus)

		// Status check endpoints
		kafkaGroup.GET("/producer/status", r.CheckProducerStatus)
		kafkaGroup.GET("/consumer/status", r.CheckConsumerStatus)
	}
}

// // ProducerRequest handles POST /kafka/producer/request
// func (r *V1) ProducerRequest(c *gin.Context) {
// 	var req struct {
// 		Topic string      `json:"topic" binding:"required"`
// 		Key   string      `json:"key"`
// 		Value interface{} `json:"value" binding:"required"`
// 	}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := r.kafka.ProduceMessage(c.Request.Context(), req.Topic, req.Key, req.Value)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(200, gin.H{"status": "message sent"})
// }

// // ConsumerReceiver handles GET /kafka/consumer/receiver
// func (r *V1) ConsumerReceiver(c *gin.Context) {
// 	topic := c.Query("topic")
// 	group := c.Query("group")
// 	if topic == "" || group == "" {
// 		c.JSON(400, gin.H{"error": "topic and group are required"})
// 		return
// 	}

// 	key, value, err := r.kafka.ConsumeMessage(c.Request.Context(), topic, group)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if value == nil {
// 		c.JSON(504, gin.H{"error": "no message received"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"key": key, "value": string(value)})
// }

// NewRedisRoutes -.
func NewRedisRoutes(apiV1Group *gin.RouterGroup, r usecase.Redis, l logger.Interface, shipperLocation usecase.ShipperLocation) {
	v1 := &V1{redis: r, l: l, v: validator.New(validator.WithRequiredStructEnabled()), shipperLocation: shipperLocation}

	redisGroup := apiV1Group.Group("/redis")
	{
		redisGroup.POST("/set", v1.setValue)
		redisGroup.GET("/get/:key", v1.getValue)
		redisGroup.POST("/shipper/location", v1.UpdateShipperLocation)
		redisGroup.GET("/shipper/location/:shipper_id", v1.GetShipperLocation)
	}
}

// NewNatsRoutes -.
func NewNatsRoutes(apiV1Group *gin.RouterGroup, nats usecase.Nats, l logger.Interface) {
	v1 := &V1{nats: nats, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	natsGroup := apiV1Group.Group("/nats")
	{
		natsGroup.POST("/publish/:subject", v1.NatsPublish)
		natsGroup.GET("/subscribe/:subject", v1.NatsSubscribe)
	}
}

// NewVietQRRoutes -.
func NewVietQRRoutes(apiV1Group *gin.RouterGroup, vietqr usecase.VietQR, l logger.Interface) {
	v1 := &V1{vietqr: vietqr, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	vietqrGroup := apiV1Group.Group("/vietqr")
	{
		vietqrGroup.POST("/gen", v1.generateQR)
		vietqrGroup.GET("/inquiry/:id", v1.inquiryQR)
		vietqrGroup.PUT("/update/:id", v1.updateStatus)
	}
}

// NewBillingRoutes -.
func NewBillingRoutes(apiV1Group *gin.RouterGroup, billing usecase.Billing, l logger.Interface) {
	// This function is deprecated - use RegisterBillingRoutes instead
	// The billing functionality is now handled by BillingController
}
