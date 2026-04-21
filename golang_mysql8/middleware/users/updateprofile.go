package middleware

import (
	"encoding/json"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

// @Summary User Profile Update
// @Description This will update user profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Param body body dto.ProfileData true "New Profile Details"
// @Success 200 {array} dto.ProfileData
// @Router /api/updateprofile/{id} [patch]
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var userDto dto.ProfileData
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	user, err := utils.GetByUserId(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if len(user) > 0 {
		db := config.Connection()
		db.Model(&models.User{}).Where("id = ?", id).Updates(userDto)
		db.Commit()

		// --- KAFKA IMPLEMENTATION ---
		kafkaProducer, err := config.GetKafkaProducer()
		if err == nil {
			defer kafkaProducer.Close()

			topic := "user-updateprofile"
			payload, _ := json.Marshal(map[string]string{
				"id": user[0].Id,
			})

			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          payload,
			}, nil)

			// Optional: Flush to ensure message is delivered
			kafkaProducer.Flush(15 * 1000)
		}

		c.JSON(200, gin.H{"message": "Your Profile has been successfully changed."})
	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}

}
