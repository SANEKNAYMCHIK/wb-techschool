package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/SANEKNAYMCHIK/order-service/internal/cache"
	"github.com/SANEKNAYMCHIK/order-service/internal/db"
	"github.com/SANEKNAYMCHIK/order-service/internal/models"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	db       *db.Postgres
	cache    *cache.LRUCache
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewConsumer(brokers []string, topic string, db *db.Postgres, cache *cache.LRUCache) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "order-service-group",
		}),
		db:       db,
		cache:    cache,
		shutdown: make(chan struct{}),
	}
}

func (c *Consumer) Start() {
	c.wg.Add(1)
	go c.run()
}

func (c *Consumer) run() {
	defer c.wg.Done()
	log.Println("Kafka consumer started")

	for {
		select {
		case <-c.shutdown:
			log.Println("Shutdown signal received")
			return
		default:
			msg, err := c.reader.FetchMessage(context.Background())
			if err != nil {
				log.Printf("Fetch error: %v", err)
				continue
			}
			c.processMessage(msg)
		}
	}
}

func (c *Consumer) processMessage(msg kafka.Message) {
	var order models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return
	}

	if err := c.db.SaveOrder(context.Background(), order); err != nil {
		log.Printf("DB save error: %v", err)
		return
	}

	c.cache.Set(order.OrderUID, order)

	if err := c.reader.CommitMessages(context.Background(), msg); err != nil {
		log.Printf("Commit error: %v", err)
	} else {
		log.Printf("Processed order %s", order.OrderUID)
	}
}

func (c *Consumer) Stop() {
	close(c.shutdown)
	c.wg.Wait()

	if err := c.reader.Close(); err != nil {
		log.Printf("Close error: %v", err)
	}
	log.Println("Kafka consumer stopped")
}
