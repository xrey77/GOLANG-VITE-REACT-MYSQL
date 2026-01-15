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

}
