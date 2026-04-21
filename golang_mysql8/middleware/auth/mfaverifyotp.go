package middleware

import (
	"encoding/json"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// @Summary MFA TOTP Verification
// @Description Multi-Factor Authenticator, OTP verification
// @Tags MultiFactor Authenticator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Param body body dto.MfaKeys true "Enter OTP Code"
// @Success 200 {array} dto.Users
// @Router /api/mfa/verifytotp/{id} [patch]
func MfaVerifyotp(c *gin.Context) {
	id := c.Param("id")

	var mfa dto.MfaKeys
	if err := c.ShouldBindJSON(&mfa); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	db := config.Connection()
	var users []dto.Users
	user := db.Where("id = ?", id).Find(&users)
	if user.Error != nil {
		c.JSON(400, gin.H{
			"message": user.Error})

		return
	}
	secret := users[0].Secret

	if len(users) > 0 {

		valid := totp.Validate(mfa.Otp, *secret)
		if valid {

			// --- KAFKA IMPLEMENTATION ---
			kafkaProducer, err := config.GetKafkaProducer()
			if err == nil {
				defer kafkaProducer.Close()

				topic := "user-verifytotp"
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

			c.JSON(200, gin.H{
				"username": users[0].Username,
				"message":  "OTP code is successfully validated."})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP code, please try again."})
			return
		}

	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}

}
