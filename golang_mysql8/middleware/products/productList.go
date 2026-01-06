package middleware

import (
	"math"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"
	"src/golang_mysql8/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Product Listings
// @Description Products Pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param page path int true "Page number"
// @Success 200 {array} []dto.Products
// @Router /products/list/:page [get]
func ProductList(c *gin.Context) {
	page := c.Param("page")
	perPage := 5
	db := config.Connection()

	var products []models.Product
	result := db.Find(&products)

	totrecs := result.RowsAffected
	total1 := float64(totrecs) / float64(perPage)
	totalPages := math.Ceil(total1)
	pg, _ := strconv.Atoi(page)
	offset := (pg - 1) * perPage

	var prods []dto.Products

	db.Limit(perPage).Offset(offset).Find(&prods)

	c.JSON(200, gin.H{
		"page":         page,
		"totpage":      totalPages,
		"totalrecords": totrecs,
		"products":     prods,
	})

}
