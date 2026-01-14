package config

import (
	"fmt"
	"log"
	"os"
	"src/golang_mysql8/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() *gorm.DB {
	host := os.Getenv("DB_HOST") // "host.docker.internal"
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, dbname)
	// fmt.Printf("DEBUG DSN: [%#v]\n", dsn1)

	// 	user, pass, host, port, dbname)
	// fmt.Printf("DEBUG DSN: [%#v]\n", dsn)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
	// 	user, pass, host, port, dbname)

	// dsn := "rey:rey@tcp(192.168.1.16:3306)/golang1.24_vitereact?charset=utf8&parseTime=True&loc=Local"

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Could not connect to MySQL database.")
	}

	fmt.Println("connected to mysql database.")
	// db.AutoMigrate(&User{}, &Role{})
	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Product{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Print("Tables Created....")

	return DB
	// var err error
	// // REMOVE the colon before '=' to use the global DB variable
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("Could not connect to MySQL database: %v", err)
	// }

	// // ... migration logic ...
	// return DB

}

// func connect() {
// 	// Replace with your actual database credentials
// 	dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("failed to connect database:", err)
// 	}

// 	fmt.Println("Successfully connected to MySQL database!")

// 	// You can now use the 'db' instance for GORM operations (e.g., AutoMigrate, Create, Find)

// 	// Ensure the connection is closed when the main function exits
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		log.Fatal("failed to get underlying sql.DB:", err)
// 	}
// 	defer sqlDB.Close()
// }
