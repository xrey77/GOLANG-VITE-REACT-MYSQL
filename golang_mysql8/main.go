package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"src/golang_mysql8/config"
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

func init() {
	err1 := godotenv.Load(".env")
	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}
	// config.Connection()
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

	producer, _ := config.GetKafkaProducer()
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.KafkaMiddleware(producer))

	router.POST("/any-endpoint", func(c *gin.Context) {
		// Retrieve producer from context
		p := c.MustGet("kafkaProducer").(*kafka.Producer)
		topic := "central-topic"

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("Your message here"),
		}, nil)

		c.JSON(200, gin.H{"status": "message sent"})
	})

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
	router.GET("/productreport", products.GetProductdata)

	router.GET("/categories", func(c *gin.Context) {
		cats, _ := products.GetAllMasterDetails()
		c.JSON(200, cats)
	})

	router.GET("/categories/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		cat, _ := products.GetCategoryWithProducts(id)
		c.JSON(200, cat)
	})

	router.GET("/chartdata", products.GetSales)

	host := "0.0.0.0"
	port := "5000"
	address := fmt.Sprintf("%s:%s", host, port)
	log.Print("Listening to ", address)
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", router))

	// if err := router.Run("127.0.0.1:5000"); err != nil {
	// 	log.Fatalf("failed to run server: %v", err)
	// }
}
