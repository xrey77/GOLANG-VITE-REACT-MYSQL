package utils

import (
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
)

func GetByUsername(username string) ([]models.User, error) {
	var users []models.User

	db := config.Connection()
	result := db.Where("username = ?", username).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetByUserId(id string) ([]dto.Users, error) {
	var users []dto.Users

	db := config.Connection()
	result := db.Where("id = ?", id).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
