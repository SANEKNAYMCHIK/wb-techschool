package cache

import (
	"context"
	"log"
	"time"

	"github.com/SANEKNAYMCHIK/order-service/internal/db"
)

func LoadCacheFromDB(ctx context.Context, cache *LRUCache, pg *db.Postgres, limit int) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	start := time.Now()
	log.Printf("Initializing cache from database")

	// Loading recent orders from db
	orders, err := pg.LoadRecentOrders(ctx, limit)
	if err != nil {
		log.Printf("Cache load failed: %v", err)
		return
	}
	// Setting cache
	count := 0
	for uid, order := range orders {
		cache.Set(uid, order)
		count++
	}
	log.Printf("Cache initialized with %d orders (took %v)", count, time.Since(start))
}
