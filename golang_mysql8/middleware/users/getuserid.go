package middleware

import (
	"encoding/json"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

// @Summary Get user by ID
// @Description Retrieve a single user's details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Success 200 {object} dto.Users
// @Router /api/getuserid/{id} [get]
func GetUserid(c *gin.Context) {
	id := c.Param("id")

	var user []dto.Users

	db := config.Connection()
	result := db.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"message": "User ID not found."})
		return
	}

	// --- KAFKA IMPLEMENTATION ---
	kafkaProducer, err := config.GetKafkaProducer()
	if err == nil {
		defer kafkaProducer.Close()

		topic := "user-getid"
		payload, _ := json.Marshal(map[string]string{
			"id":        user[0].Id,
			"firstname": user[0].Firstname,
		})

		kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          payload,
		}, nil)

		// Optional: Flush to ensure message is delivered
		kafkaProducer.Flush(15 * 1000)
	}

	c.JSON(200, user)
}
