package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SANEKNAYMCHIK/order-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, connStr string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("incorrect configuration string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool failed: %w", err)
	}

	// check connecting
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	// Create tables in database
	if err := CreateSchema(context.Background(), pool); err != nil {
		return nil, fmt.Errorf("schema creation failed: %w", err)
	}

	log.Printf("PostgreSQL connected")
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) SaveOrder(ctx context.Context, order models.Order) error {
	// Timeout for saving operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)

	// Save Order
	if err := saveOrderMain(ctx, tx, order); err != nil {
		return err
	}

	// Save Delivery
	if err := saveDelivery(ctx, tx, order); err != nil {
		return err
	}

	// Save Payments
	if err := savePayment(ctx, tx, order); err != nil {
		return err
	}

	// Save Items
	if err := saveItems(ctx, tx, order); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	log.Printf("📦 Order %s saved", order.OrderUID)
	return nil
}

func saveOrderMain(ctx context.Context, tx pgx.Tx, order models.Order) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO UPDATE SET
            track_number = EXCLUDED.track_number,
            entry = EXCLUDED.entry,
            locale = EXCLUDED.locale,
            internal_signature = EXCLUDED.internal_signature,
            customer_id = EXCLUDED.customer_id,
            delivery_service = EXCLUDED.delivery_service,
            shardkey = EXCLUDED.shardkey,
            sm_id = EXCLUDED.sm_id,
            date_created = EXCLUDED.date_created,
            oof_shard = EXCLUDED.oof_shard`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)

	if err != nil {
		return fmt.Errorf("save order failed: %w", err)
	}
	return nil
}

func saveDelivery(ctx context.Context, tx pgx.Tx, order models.Order) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (order_uid) DO UPDATE SET
            name = EXCLUDED.name,
            phone = EXCLUDED.phone,
            zip = EXCLUDED.zip,
            city = EXCLUDED.city,
            address = EXCLUDED.address,
            region = EXCLUDED.region,
            email = EXCLUDED.email`,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	if err != nil {
		return fmt.Errorf("save delivery failed: %w", err)
	}
	return nil
}

func savePayment(ctx context.Context, tx pgx.Tx, order models.Order) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO payments (
            transaction, order_uid, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (transaction) DO UPDATE SET
            order_uid = EXCLUDED.order_uid,
            request_id = EXCLUDED.request_id,
            currency = EXCLUDED.currency,
            provider = EXCLUDED.provider,
            amount = EXCLUDED.amount,
            payment_dt = EXCLUDED.payment_dt,
            bank = EXCLUDED.bank,
            delivery_cost = EXCLUDED.delivery_cost,
            goods_total = EXCLUDED.goods_total,
            custom_fee = EXCLUDED.custom_fee`,
		order.Payment.Transaction,
		order.OrderUID,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)

	if err != nil {
		return fmt.Errorf("save payment failed: %w", err)
	}
	return nil
}

func saveItems(ctx context.Context, tx pgx.Tx, order models.Order) error {
	batch := &pgx.Batch{}

	for _, item := range order.Items {
		batch.Queue(`
			INSERT INTO items (
                chrt_id, order_uid, track_number, price, rid, 
                name, sale, size, total_price, nm_id, brand, status
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            ON CONFLICT (chrt_id) DO UPDATE SET
                track_number = EXCLUDED.track_number,
                price = EXCLUDED.price,
                rid = EXCLUDED.rid,
                name = EXCLUDED.name,
                sale = EXCLUDED.sale,
                size = EXCLUDED.size,
                total_price = EXCLUDED.total_price,
                nm_id = EXCLUDED.nm_id,
                brand = EXCLUDED.brand,
                status = EXCLUDED.status`,
			item.ChrtID,
			order.OrderUID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range order.Items {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("save item failed: %w", err)
		}
	}

	return nil
}

func (p *Postgres) LoadRecentOrders(ctx context.Context, limit int) (map[string]models.Order, error) {
	// 1. request for main data of order
	orderQuery := `
        SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
               o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, 
               p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        JOIN deliveries d ON o.order_uid = d.order_uid
        JOIN payments p ON o.order_uid = p.order_uid
        ORDER BY o.date_created DESC
        LIMIT $1`

	rows, err := p.pool.Query(ctx, orderQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent orders: %w", err)
	}
	defer rows.Close()

	orders := make(map[string]models.Order)
	orderUIDs := make([]string, 0, limit)

	// 2. Adding main data of Order
	for rows.Next() {
		var o models.Order
		o.Items = make([]models.Item, 0)
		if err := rows.Scan(
			&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
			&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
			&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
			&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency,
			&o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank,
			&o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders[o.OrderUID] = o
		orderUIDs = append(orderUIDs, o.OrderUID)
	}

	// 3. Download Items
	if len(orderUIDs) > 0 {
		itemsQuery := `
            SELECT order_uid, chrt_id, track_number, price, rid, 
                   name, sale, size, total_price, nm_id, brand, status
            FROM items
            WHERE order_uid = ANY($1)`

		itemRows, err := p.pool.Query(ctx, itemsQuery, orderUIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to query items: %w", err)
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var i models.Item
			var orderUID string
			if err := itemRows.Scan(
				&orderUID, &i.ChrtID, &i.TrackNumber, &i.Price, &i.RID,
				&i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NmID, &i.Brand, &i.Status,
			); err != nil {
				return nil, fmt.Errorf("failed to scan item: %w", err)
			}

			if order, exists := orders[orderUID]; exists {
				order.Items = append(order.Items, i)
				orders[orderUID] = order
			}
		}
	}

	return orders, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
	log.Println("PostgreSQL connection closed")
}

func (p *Postgres) GetOrderByUID(ctx context.Context, uid string) (*models.Order, error) {
	query := `
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
               o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
               d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
               p.transaction, p.request_id, p.currency, p.provider, p.amount, 
               p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        JOIN deliveries d ON o.order_uid = d.order_uid
        JOIN payments p ON o.order_uid = p.order_uid
		WHERE o.order_uid = $1`

	var order models.Order
	err := p.pool.QueryRow(ctx, query, uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank,
		&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}
	itemsQuery := `
	SELECT order_uid, chrt_id, track_number, price, rid, 
	name, sale, size, total_price, nm_id, brand, status
	FROM items
	WHERE order_uid = $1`
	rows, err := p.pool.Query(ctx, itemsQuery, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	order.Items = items
	return &order, nil
}
