package middleware

import (
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/gin-gonic/gin"
)

// @Summary Get user by ID
// @Description Retrieve a single user's details
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Success 200 {object} dto.Users
// @Router /api/getuserid/{id} [get]
func GetUserid(c *gin.Context) {
	id := c.Param("id")

	var user []dto.Users

	db := config.Connection()
	result := db.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		// GORM returns ErrRecordNotFound if nothing is found
		c.JSON(404, gin.H{"message": "User ID not found."})
		return
	}

	c.JSON(200, user)
}
