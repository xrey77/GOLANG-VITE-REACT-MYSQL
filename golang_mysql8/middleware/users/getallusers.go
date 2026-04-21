package middleware

import (
	"encoding/json"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

// @Summary Retrieve users
// @Description Display all users
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Users
// @Router /api/getallusers [get]
func GetAllusers(c *gin.Context) {

	db := config.Connection()
	var users []dto.Users

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve records"})
		return
	}

	// --- KAFKA IMPLEMENTATION ---
	kafkaProducer, err := config.GetKafkaProducer()
	if err == nil {
		defer kafkaProducer.Close()

		topic := "users-getall"
		payload, _ := json.Marshal(map[string]string{
			"id": users[0].Id,
		})

		kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          payload,
		}, nil)

		// Optional: Flush to ensure message is delivered
		kafkaProducer.Flush(15 * 1000)
	}

	c.JSON(200, users)

}
