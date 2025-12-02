package middleware

import (
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
)

func GetUserid(c *gin.Context) {
	id := c.Param("id")
	user, err := utils.GetByUserId(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if len(user) == 0 {
		c.JSON(400, gin.H{"message": "User ID not found."})
		return
	}

	c.JSON(200, user)
}
