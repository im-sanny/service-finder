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
		RETURNING id`, // RETURNING id gives you one newly-created ID.
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
	// The loop basically means:
	// "Give me the first row → scan it → put it in my slice.
	// Give me the next row → scan it → put it in my slice.
	// Keep going until there are no more rows."
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

func ServicePatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	SET name=COALESCE($1, name),
	description=COALESCE($2, description)
	WHERE id=$3
	RETURNING id, name, description`, // why i need to return these for patch when i don't need it for update?
		s.Name, s.Description, id).Scan(&s.ID, &s.Name, &s.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to patch service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func ServiceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var deletedId int
	err = db.QueryRow(`DELETE FROM services WHERE id=$1 RETURNING id`, id).Scan(&deletedId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func ProviderPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	err := db.QueryRow(`INSERT INTO providers (name, phone, location, description, service_id) VALUES($1, $2, $3, $4, $5) RETURNING id`, p.Name, p.Phone, p.Location, p.Description, p.ServiceID).Scan(&p.ID)

	if err != nil {
		http.Error(w, "Failed to insert service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
func ProviderGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`SELECT name, phone, location, description, service_id  FROM providers`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var provider []Provider
	for rows.Next() {
		var p Provider
		rows.Scan(&p.Name, &p.Phone, &p.Description, &p.ServiceID)
		provider = append(provider, p)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)

}
func ProviderID(w http.ResponseWriter, r *http.Request) {
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

	var p Provider
	err = db.QueryRow(`SELECT id, name, phone, location, description, service_id FROM providers WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Phone, &p.Description, &p.ServiceID)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
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
	mux.HandleFunc("PATCH /service/{id}", ServicePatch)
	mux.HandleFunc("DELETE /service/{id}", ServiceDelete)

	mux.HandleFunc("POST /provider", ProviderPost)
	mux.HandleFunc("GET /provider", ProviderGet)
	mux.HandleFunc("GET /provider/{id}", ProviderID)

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
