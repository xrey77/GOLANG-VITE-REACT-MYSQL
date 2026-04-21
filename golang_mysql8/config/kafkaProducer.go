package config

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func GetKafkaProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "registration-service",
		"acks":              "all",
		"security.protocol": "SASL_PLAIN",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "guests",
		"sasl.password":     "guests",
	})
	if err != nil {
		return nil, err
	}

	// Essential: Drain events in background
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()
	return p, nil
}

// func GetKafkaProducer() (*kafka.Producer, error) {
// 	return kafka.NewProducer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"client.id":         "registration-service",
// 		"acks":              "all",

// 		// Authentication Settings
// 		"security.protocol": "SASL_SSL", // "SASL_PLAINTEXT" if not using SSL
// 		"sasl.mechanism":    "PLAIN",    // Common mechanisms: "PLAIN", "SCRAM-SHA-256", or "SCRAM-SHA-512"
// 		"sasl.username":     "guests",
// 		"sasl.password":     "guests",
// 	})
// }
