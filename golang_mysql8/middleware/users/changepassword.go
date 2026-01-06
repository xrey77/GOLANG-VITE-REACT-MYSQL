package middleware

import (
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
)

// @Summary Change User Password
// @Description User Change Password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Param body body dto.ChangePassword true "New Password Details"
// @Success 200 {object} dto.ChangePassword
// @Router /api/changepassword/{id} [patch]
func ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var userDto dto.ChangePassword

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
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
