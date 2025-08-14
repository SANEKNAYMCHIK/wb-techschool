package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create all tables with initialization
func CreateSchema(ctx context.Context, conn *pgxpool.Pool) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			order_uid          TEXT PRIMARY KEY,
			track_number       TEXT NOT NULL,
			entry              TEXT NOT NULL,
			locale             TEXT NOT NULL,
			internal_signature TEXT,
			customer_id        TEXT NOT NULL,
			delivery_service   TEXT NOT NULL,
			shardkey           TEXT NOT NULL,
			sm_id              INTEGER NOT NULL,
			date_created       TIMESTAMPTZ NOT NULL,
			oof_shard          TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create orders table failed: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS deliveries (
			order_uid   TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
			name        TEXT NOT NULL,
			phone       TEXT NOT NULL,
			zip         TEXT NOT NULL,
			city        TEXT NOT NULL,
			address     TEXT NOT NULL,
			region      TEXT NOT NULL,
			email       TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create deliveries table failed: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS payments (
			transaction   TEXT PRIMARY KEY,
			order_uid     TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
			request_id    TEXT,
			currency      TEXT NOT NULL,
			provider      TEXT NOT NULL,
			amount        INTEGER NOT NULL,
			payment_dt    BIGINT NOT NULL,
			bank          TEXT NOT NULL,
			delivery_cost INTEGER NOT NULL,
			goods_total   INTEGER NOT NULL,
			custom_fee    INTEGER NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create payments table failed: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS items (
			chrt_id      BIGINT PRIMARY KEY,
			order_uid    TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
			track_number TEXT NOT NULL,
			price        INTEGER NOT NULL,
			rid          TEXT NOT NULL,
			name         TEXT NOT NULL,
			sale         INTEGER NOT NULL,
			size         TEXT NOT NULL,
			total_price  INTEGER NOT NULL,
			nm_id        BIGINT NOT NULL,
			brand        TEXT NOT NULL,
			status       INTEGER NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create items table failed: %w", err)
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_payment_order ON payments (order_uid);
		CREATE INDEX IF NOT EXISTS idx_items_order ON items (order_uid);
	`)
	if err != nil {
		return fmt.Errorf("create indexes failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	log.Println("Database schema created")
	return nil
}
