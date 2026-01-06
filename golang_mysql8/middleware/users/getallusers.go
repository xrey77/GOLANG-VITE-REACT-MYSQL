package middleware

import (
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Retrieve users
// @Description Display all users
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.Users
// @Router /api/getallusers [get]
func GetAllusers(c *gin.Context) {

	db := config.Connection()
	var users []dto.Users

	// GORM's Find method to retrieve all records
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve records"})
		return
	}

	c.JSON(200, users)

}
