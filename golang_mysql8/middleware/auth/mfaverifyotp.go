package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func MfaVerifyotp(c *gin.Context) {
	id := c.Param("id")

	var mfa dto.MfaKeys
	err1 := json.NewDecoder(c.Request.Body).Decode(&mfa)

	if err1 != nil {
		log.Fatalf("Unable to decode the request body.  %v", err1)
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
			c.JSON(200, gin.H{
				"username": users[0].Username,
				"message":  "OTP code is successfully validated.s"})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP code, please try again."})
			return
		}

	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}

	// var req OTPRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// totp.Validate checks if the provided OTP is currently valid within the time step (usually 30s)

}
