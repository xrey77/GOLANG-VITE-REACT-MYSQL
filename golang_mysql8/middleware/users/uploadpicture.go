package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"src/golang_mysql8/config"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
)

// @Summary Update user profile picture
// @Description Upload user picture
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "User Id"
// @Param userpic formData file true "New Profile Picture"
// @Success 200 {object} map[string]interface{}
// @Router /api/uploadpicture/{id} [patch]
func UploadPicture(c *gin.Context) {
	id := c.Param("id")
	user, err := utils.GetByUserId(id)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	if len(user) > 0 {

		file, err := c.FormFile("userpic") // "file" is the key for the form data
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		newfile := "00" + id + ext
		dst := filepath.Join("./assets/users/", newfile) // Destination path

		// Save the uploaded file to the specified destination
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		db := config.Connection()
		uri := "http://localhost:5000/assets/users/" + newfile
		db.Model(&models.User{}).Where("id = ?", id).Update("userpicture", uri)
		db.Commit()

		c.JSON(200, gin.H{
			"userpic": uri,
			"message": "Profile picture has been changed."})

	} else {
		c.JSON(400, gin.H{"message": "User ID not found."})
	}

}
