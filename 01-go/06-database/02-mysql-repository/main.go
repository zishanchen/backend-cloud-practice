package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewMySQLTaskRepository(db)
	handler := NewTaskHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /tasks",
		handler.List,
	)

	mux.HandleFunc(
		"GET /tasks/{id}",
		handler.Get,
	)

	mux.HandleFunc(
		"POST /tasks",
		handler.Create,
	)

	server := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println(
		"MySQL task API running on http://localhost:8081",
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")

	db, err := sql.Open(
		"mysql",
		dsn,
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(
		3 * time.Minute,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
