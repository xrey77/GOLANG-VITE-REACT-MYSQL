package middleware

import (
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
)

func GetCategoryWithProducts(catID int) (dto.Category, error) {
	var category dto.Category
	db := config.Connection()

	err := db.Preload("Products").First(&category, catID).Error

	return category, err
}

func GetAllMasterDetails() ([]dto.Category, error) {
	var categories []dto.Category
	db := config.Connection()

	err := db.Preload("Products").Find(&categories).Error

	return categories, err
}
