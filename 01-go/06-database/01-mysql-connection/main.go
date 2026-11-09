package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	db, err := openDB(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	fmt.Println("MySQL connection OK")

	printPoolStats(db)
}

func openDB(ctx context.Context) (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, fmt.Errorf(
			"MYSQL_DSN environment variable is not set",
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf(
			"open database handle: %w",
			err,
		)
	}

	configurePool(db)

	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf(
			"ping database: %w",
			err,
		)
	}

	return db, nil
}

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	db.SetConnMaxLifetime(
		3 * time.Minute,
	)

	db.SetConnMaxIdleTime(
		1 * time.Minute,
	)
}

func printPoolStats(db *sql.DB) {
	stats := db.Stats()

	fmt.Println("Connection pool:")
	fmt.Printf(
		"  max open: %d\n",
		stats.MaxOpenConnections,
	)
	fmt.Printf(
		"  open: %d\n",
		stats.OpenConnections,
	)
	fmt.Printf(
		"  in use: %d\n",
		stats.InUse,
	)
	fmt.Printf(
		"  idle: %d\n",
		stats.Idle,
	)
}
