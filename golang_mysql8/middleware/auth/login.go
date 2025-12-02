package middleware

import (
	"encoding/json"
	"errors"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var userDto dto.UserLogin
	err := json.NewDecoder(c.Request.Body).Decode(&userDto)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "Unable to decode the request body."})
	}
	plainPwd := userDto.Password
	user, err := GetUserInfo(userDto.Username)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error()})
		return
	}
	if user.Username != "" {
		hashPwd := user.Password
		err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plainPwd))
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid Password."})
		} else {

			token, _ := utils.GenerateJWT(user.Email)

			c.JSON(200, gin.H{
				"id":          user.Id,
				"firstname":   user.Firstname,
				"lastname":    user.Lastname,
				"email":       user.Email,
				"mobile":      user.Mobile,
				"username":    user.Username,
				"roles":       user.Roles,
				"isactivated": user.Isactivated,
				"isblocked":   user.Isblocked,
				"userpic":     user.Userpicture,
				"qrcodeurl":   user.Qrcodeurl,
				"token":       token,
				"message":     "Login Successfull."})
		}

	} else {
		c.JSON(200, gin.H{
			"message": "Username not found, please register."})
	}
}

func GetUserInfo(userName string) (*dto.Users, error) {
	var user dto.Users

	db := config.Connection()
	result := db.Where("username = ?", userName).Find(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, result.Error
	}
	return &user, nil
}
