package middleware

import (
	"errors"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary User Login
// @Description Authenticat User
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.UserLogin true "User Login Credentials"
// @Success 200 {array} dto.UserLogin
// @Router /auth/signin [post]
func Login(c *gin.Context) {
	var userDto dto.UserLogin
	// err := json.NewDecoder(c.Request.Body).Decode(&userDto)

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	plainPwd := userDto.Password
	user, err := GetUserInfo(userDto.Username)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if user.Username != "" {

		hashPwd := user.Password
		err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plainPwd))
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid Password."})
			return
		} else {

			token, _ := utils.GenerateJWT(user.Email)
			roleDto, _ := GetRolName(user.Id)
			var rolesName string
			if roleDto != nil {
				rolesName = roleDto.Name
			}

			c.JSON(200, gin.H{
				"id":          user.Id,
				"firstname":   user.Firstname,
				"lastname":    user.Lastname,
				"email":       user.Email,
				"mobile":      user.Mobile,
				"username":    user.Username,
				"roles":       rolesName,
				"isactivated": user.Isactivated,
				"isblocked":   user.Isblocked,
				"userpic":     user.Userpicture,
				"qrcodeurl":   user.Qrcodeurl,
				"token":       token,
				"message":     "Login Successfull."})
		}

	} else {
		c.JSON(404, gin.H{
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

func GetRolName(id string) (*dto.Roles, error) {
	var roles dto.Roles

	db := config.Connection()
	result := db.Where("id = ?", id).Find(&roles)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Role not found
		}
		return nil, result.Error
	}
	return &roles, nil
}
