package middleware

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

func KafkaMiddleware(producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("kafkaProducer", producer)
		c.Next()
	}
}
