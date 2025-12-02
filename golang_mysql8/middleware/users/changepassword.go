package middleware

import (
	"encoding/json"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
)

func ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var userDto dto.ChangePassword
	err := json.NewDecoder(c.Request.Body).Decode(&userDto)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "Unable to decode the request body."})
	}

	user, err := utils.GetByUserId(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if len(user) > 0 {
		hash, _ := utils.HashPassword(userDto.Password)
		db := config.Connection()
		db.Model(&models.User{}).Where("id = ?", id).Update("password", hash)
		db.Commit()
		c.JSON(200, gin.H{"message": "Password has been changed."})
	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}
}
