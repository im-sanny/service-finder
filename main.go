package main

import (
	"database/sql"
	"encoding/json"
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Phone       int    `json:"phone"`
	Location    string `json:"location"`
	Description string `json:"description"`
	ServiceID   int    `json:"service_id"`
}

var db *sql.DB

func ServicePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var s Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(`
		INSERT INTO services (name, description)
		VALUES ($1, $2)
		RETURNING id`,
		s.Name, s.Description).Scan(&s.ID)

	if err != nil {
		http.Error(w, "Failed to insert service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

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

	mux.HandleFunc("POST /service", ServicePost)

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
