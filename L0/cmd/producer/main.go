package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SANEKNAYMCHIK/order-service/internal/models"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")

	// Creating Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
	})

	defer writer.Close()

	log.Println("Producer started. Press Ctrl+C to stop.")

	// Every 3 seconds send an order
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down producer...")
			return
		case <-ticker.C:
			// Generating test orders
			orderGenerator := models.NewTestOrderGenerator()
			order := orderGenerator.GenerateOrder()
			value, err := json.Marshal(order)
			if err != nil {
				log.Printf("Error marshaling order: %v", err)
				continue
			}

			err = writer.WriteMessages(ctx, kafka.Message{
				Value: value,
			})
			if err != nil {
				log.Printf("Failed to write message: %v", err)
			} else {
				log.Printf("Sent order %s", order.OrderUID)
			}
		}
	}
}
