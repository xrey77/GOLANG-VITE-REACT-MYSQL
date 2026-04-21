package middleware

import (
	"encoding/json"
	"net/http"
	"src/golang_mysql8/config"
	"src/golang_mysql8/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
)

// @Summary Retrieve sales
// @Description Display all sales
// @Tags Sale
// @Produce json
// @Success 200 {array} dto.Sales
// @Router /api/chartdata [get]
func GetSales(c *gin.Context) {

	db := config.Connection()
	var sales []dto.Sales

	if err := db.Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve records"})
		return
	}

	// --- KAFKA IMPLEMENTATION ---
	kafkaProducer, err := config.GetKafkaProducer()
	if err == nil {
		defer kafkaProducer.Close()

		topic := "sales-getall"
		payload, _ := json.Marshal(map[string]string{
			"id": sales[0].Id,
		})

		kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          payload,
		}, nil)

		// Optional: Flush to ensure message is delivered
		kafkaProducer.Flush(15 * 1000)
	}
	c.JSON(200, sales)
}
