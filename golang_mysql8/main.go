package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "src/golang_mysql8/docs"
	"src/golang_mysql8/middleware"
	auth "src/golang_mysql8/middleware/auth"
	products "src/golang_mysql8/middleware/products"
	users "src/golang_mysql8/middleware/users"

	"github.com/gin-gonic/contrib/static"

	swaggerFiles "github.com/swaggo/files"
	// Add the closing quote and full path below
	_ "src/golang_mysql8/docs" // Side-effect import for generated docs

	ginSwagger "github.com/swaggo/gin-swagger"
)

// Add this line

func init() {
	// config.Connection()
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

}

// @title BARCLAYS BANK API Management
// @version 1.0
// @description REST API Documentation Gin server.
// @host localhost:5000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(static.Serve("/jesuskingofkings", static.LocalFile("templates", true)))
	router.Static("/assets", "./assets")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost", "http://localhost:5000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/auth/signin", auth.Login)
	router.POST("/auth/signup", auth.Register)

	authGuard := router.Group("/api")
	authGuard.Use(middleware.AuthMiddleware())
	{
		authGuard.GET("/getallusers", users.GetAllusers)
		authGuard.GET("/getuserid/:id", users.GetUserid)
		authGuard.PATCH("/changepassword/:id", users.ChangePassword)
		authGuard.PATCH("/updateprofile/:id", users.UpdateProfile)
		authGuard.PATCH("/uploadpicture/:id", users.UploadPicture)
		authGuard.PATCH("/mfa/activate/:id", auth.MfaActivate)
		authGuard.PATCH("/mfa/verifytotp/:id", auth.MfaVerifyotp)
	}

	router.GET("/products/list/:page", products.ProductList)
	router.GET("/products/search/:page/:key", products.ProductSearch)

	host := "0.0.0.0"
	port := "5000"
	address := fmt.Sprintf("%s:%s", host, port)
	log.Print("Listening to ", address)
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", router))

	// if err := router.Run("127.0.0.1:5000"); err != nil {
	// 	log.Fatalf("failed to run server: %v", err)
	// }
}
