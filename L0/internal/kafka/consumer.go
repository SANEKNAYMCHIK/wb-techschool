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
	c.createTopic()
	c.wg.Add(1)
	go c.run()
}

func (c *Consumer) createTopic() {
	log.Print(c.reader.Config().Brokers[0])
	log.Print(c.reader.Config().Topic)
	// 1. Connecting to Kafka
	conn, err := kafka.Dial("tcp", c.reader.Config().Brokers[0])
	if err != nil {
		log.Printf("Failed to connect to Kafka: %v", err)
		return
	}
	defer conn.Close()

	// 2. Check existing of the topic
	partitions, err := conn.ReadPartitions(c.reader.Config().Topic)
	if err == nil && len(partitions) > 0 {
		log.Printf("Topic '%s' already exists", c.reader.Config().Topic)
		return
	}

	// 3. Create topic
	log.Printf("Creating topic '%s'...", c.reader.Config().Topic)
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             c.reader.Config().Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	if err != nil {
		log.Printf("Failed to create topic: %v", err)
	} else {
		log.Printf("Topic successfully created")
	}
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
	log.Print("Get a new message: ")
	log.Print(msg)
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
