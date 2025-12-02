package middleware

import (
	"encoding/json"
	"log"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user dto.UserRegister
	err := json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	plainPwd := user.Password
	hashPwd, _ := utils.HashPassword(plainPwd)

	db := config.Connection()
	userEmail, _ := SearchByEmail(user.Email)
	if len(userEmail) > 0 {
		c.JSON(400, gin.H{
			"message": "Email Address is already taken."})
		return
	}

	userName, _ := SearchByUsername(user.Username)
	if len(userName) > 0 {
		c.JSON(400, gin.H{"message": "Username is already taken."})
		return
	}

	userModel := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Username:  user.Username,
		Password:  hashPwd,
	}
	res := db.Create(&userModel)
	rec := res.RowsAffected
	if rec > 0 {
		c.JSON(200, gin.H{"message": "Registration Successfull, please login now."})

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
