package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	utils "src/golang_mysql8/util"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

func MfaActivate(c *gin.Context) {
	id := c.Param("id")
	var user dto.MfaActivation
	err := json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	db := config.Connection()

	if user.TwoFactoEnabled {
		user, err := utils.GetByUserId(id)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if len(user) > 0 {
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "BARCLAYS BANK", // The name of your application
				AccountName: user[0].Email,   // The user's account identifier
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
				return
			}
			// The key.Secret() is the base32 encoded secret you must save
			secret := key.Secret()
			// The key.URL() is the otpauth URI, which can be converted into a QR code
			qrCodeURL := key.URL()

			pngBytes, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate QR code: %v", err)})
				return
			}
			// 3. Base64 encode the PNG bytes
			var mfaData dto.MfaData
			// "data:image/png;base64,
			base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
			mfaData.Secret = secret
			mfaData.Qrcodeurl = string(base64Encoded)

			db.Model(&models.User{}).Where("id = ?", id).Updates(mfaData)
			db.Commit()

			c.JSON(200, gin.H{
				"qrcodeurl": base64Encoded,
				"message":   "Multi-Factor Authenticator has been enabled."})

		}

	} else {

		db.Model(&models.User{}).Where("id = ?", id).Update("qrcodeurl", nil)
		db.Commit()

		c.JSON(200, gin.H{
			"message": "Multi-Factor Authenticator has been disabled."})

	}

}
