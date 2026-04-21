package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

// @Summary MFA Activation
// @Description Multi-Factor Authenticator
// @Tags MultiFactor Authenticator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Param body body dto.MfaActivation true "Enable MFA"
// @Success 200 {array} dto.MfaActivation
// @Router /api/mfa/activate/{id} [patch]
func MfaActivate(c *gin.Context) {
	id := c.Param("id")
	var user dto.MfaActivation
	err := json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	db := config.Connection()

	if user.TwoFactorEnabled {
		user, err := utils.GetByUserId(id)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if len(user) > 0 {
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "BARCLAYS BANK",
				AccountName: user[0].Email,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
				return
			}

			secret := key.Secret()
			qrCodeURL := key.URL()

			pngBytes, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate QR code: %v", err)})
				return
			}

			var mfaData dto.MfaData
			base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
			mfaData.Secret = secret
			mfaData.Qrcodeurl = string(base64Encoded)

			db.Model(&models.User{}).Where("id = ?", id).Updates(mfaData)
			db.Commit()

			// --- KAFKA IMPLEMENTATION ---
			kafkaProducer, err := config.GetKafkaProducer()
			if err == nil {
				defer kafkaProducer.Close()

				topic := "user-activatemfa"
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

			c.JSON(200, gin.H{
				"qrcodeurl": base64Encoded,
				"message":   "Multi-Factor Authenticator has been enabled."})

		}

	} else {

		// --- KAFKA IMPLEMENTATION ---
		kafkaProducer, err := config.GetKafkaProducer()
		if err == nil {
			defer kafkaProducer.Close()

			topic := "user-activatemfa"
			payload, _ := json.Marshal(map[string]string{
				"id": id,
			})

			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          payload,
			}, nil)

			// Optional: Flush to ensure message is delivered
			kafkaProducer.Flush(15 * 1000)
		}

		db.Model(&models.User{}).Where("id = ?", id).Update("qrcodeurl", nil)
		db.Commit()

		c.JSON(200, gin.H{
			"message": "Multi-Factor Authenticator has been disabled."})

	}

}
