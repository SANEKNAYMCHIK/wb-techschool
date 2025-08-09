package main

import (
	"context"
	"log"

	"github.com/SANEKNAYMCHIK/order-service/internal/db"
)

func main() {
	connStr := ""

	pg, err := db.NewPostgres(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pg.Close(context.Background())

	// cache := cache.NewCache()

}
