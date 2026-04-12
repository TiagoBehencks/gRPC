package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	ctx := context.Background()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	Pool = pool
	return nil
}

func InitSchema() error {
	ctx := context.Background()
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		quantity INT NOT NULL
	);
	`
	_, err := Pool.Exec(ctx, schema)
	return err
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}