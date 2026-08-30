package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Service struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type Provider struct {
	ID          int     `json:"id"`
	Name        *string `json:"name"`
	Phone       *int    `json:"phone"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
	ServiceID   int     `json:"service_id"`
}

var db *sql.DB

func main() {
	var err error
	conStr := "postgres://postgres:360420@localhost:5432/serfin?sslmode=disable"

	db, err = sql.Open("postgres", conStr)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connected successfully!")

	mux := http.NewServeMux()

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
