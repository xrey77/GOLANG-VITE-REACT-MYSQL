package middleware

import (
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/gin-gonic/gin"
)

func GetAllusers(c *gin.Context) {

	db := config.Connection()
	var users []dto.Users

	// Use GORM's Find method to retrieve all records
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve records"})
		return
	}

	c.JSON(200, users)

}
