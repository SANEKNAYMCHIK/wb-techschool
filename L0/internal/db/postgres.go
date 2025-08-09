package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	log.Println("Connected to PostgreSQL")
	return &Postgres{conn: conn}, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}
