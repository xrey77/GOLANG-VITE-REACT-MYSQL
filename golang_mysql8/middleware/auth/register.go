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

// @Summary User Registration
// @Description Create User Account
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.UserRegister true "Account Registration"
// @Success 200 {array} dto.UserRegister
// @Router /auth/signup [post]
func Register(c *gin.Context) {
	var user dto.UserRegister
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	userEmail, err := SearchByEmail(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"message": "Database lookup failed"})
		return
	}
	if len(userEmail) > 0 {
		c.JSON(400, gin.H{
			"message": "Email Address is already taken."})
		return
	}

	userName, err := SearchByUsername(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"message": "Database lookup failed"})
		return
	}
	if len(userName) > 0 {
		c.JSON(400, gin.H{"message": "Username is already taken."})
		return
	}

	plainPwd := user.Password
	hashPwd, _ := utils.HashPassword(plainPwd)

	userModel := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Username:  user.Username,
		Password:  hashPwd,
		Role_id:   2,
	}

	db := config.Connection()
	if res := db.Create(&userModel); res.RowsAffected > 0 {

		// --- KAFKA IMPLEMENTATION ---
		kafkaProducer, err := config.GetKafkaProducer()
		if err == nil {
			defer kafkaProducer.Close()

			topic := "user-registrations"
			payload, _ := json.Marshal(map[string]string{
				"email":     user.Email,
				"firstname": user.Firstname,
			})

			kafkaProducer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          payload,
			}, nil)

			// Optional: Flush to ensure message is delivered
			kafkaProducer.Flush(15 * 1000)
		}

		c.JSON(200, gin.H{"message": "Registration Successful, please login now."})
	} else {
		c.JSON(400, gin.H{"message": "Registration Failed.."})
	}
}

func SearchByEmail(email string) ([]models.User, error) {
	var users []models.User

	db := config.Connection()
	result := db.Where("email = ?", email).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func SearchByUsername(username string) ([]models.User, error) {
	var users []models.User

	db := config.Connection()
	result := db.Where("username = ?", username).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
