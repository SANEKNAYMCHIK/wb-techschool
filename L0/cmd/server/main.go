package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/SANEKNAYMCHIK/order-service/internal/cache"
	"github.com/SANEKNAYMCHIK/order-service/internal/db"
	"github.com/SANEKNAYMCHIK/order-service/internal/httpfiles"
	"github.com/SANEKNAYMCHIK/order-service/internal/kafka"
)

func main() {
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	connStr := os.Getenv("CONN_STR")
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	cacheCapacity, _ := strconv.Atoi(os.Getenv("CACHE_CAPACITY"))

	// Init db
	pg, err := db.NewPostgres(mainCtx, connStr)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer pg.Close()

	// Init cache
	lruCache := cache.NewLRUCache(cacheCapacity)
	log.Print("Cache created")
	go cache.LoadCacheFromDB(mainCtx, lruCache, pg, cacheCapacity)

	// Kafka Consumer
	kafkaConsumer := kafka.NewConsumer(kafkaBrokers, topic, pg, lruCache)
	kafkaConsumer.Start()

	httpServer := httpfiles.NewServer(lruCache, pg)
	server := &http.Server{
		Addr:    ":8080",
		Handler: httpServer,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("HTTP server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
			serverErrors <- err
		}
	}()

	// Graceful Shutdown

	// Waiting signals
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-shutdownSignal:
		log.Printf("Received signal: %v. Starting shutdown...", sig)
	case err := <-serverErrors:
		log.Printf("Server error: %v. Starting shutdown...", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	// Using WaitGroup for stopping all components
	var wg sync.WaitGroup

	// 1. Stop HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Shutting down HTTP server...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server stopped gracefully")
		}
	}()

	// 2. Stop Kafka Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Shutting down Kafka consumer...")
		kafkaConsumer.Stop()
		log.Println("Kafka consumer stopped")
	}()

	// 3. Stop main context
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainCancel()
		log.Println("Main context canceled")
	}()

	// Waiting for stopping all operations or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All components stopped successfully")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timed out, forcing exit")
	}

	log.Println("Shutdown completed")
}
