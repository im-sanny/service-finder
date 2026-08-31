package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
		s.Name, s.Description).Scan(&s.ID) // The & means you're giving Scan the memory addresses where it should put the values.

	if err != nil {
		http.Error(w, "Failed to insert service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func ServiceGet(w http.ResponseWriter, r *http.Request) { // r request for data and w writes or provide that data
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT id, name, description FROM services;`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // why?

	services := make([]Service, 0)
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Description); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		services = append(services, s)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func ServiceId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var s Service
	err = db.QueryRow(`SELECT id, name, description, FROM services WHERE id=$1;`, id).Scan(&s.ID, &s.Name, &s.Description)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func ServicePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var s Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(`
	UPDATE services
	SET name=$1, description=$2, WHERE id=$3`,
		s.Name, s.Description, id).Scan(&s.ID, &s.Name, &s.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
	mux.HandleFunc("GET /service", ServiceGet)
	mux.HandleFunc("GET /service/{id}", ServiceId)
	mux.HandleFunc("PUT /service/{id}", ServicePut)

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
