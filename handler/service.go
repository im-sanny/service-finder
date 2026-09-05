package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/im-sanny/service-finder/model"
)

type ServiceHandler struct {
	DB *sql.DB // why use pointer for sql.DB?
}

func (h *ServiceHandler) ServicePost(w http.ResponseWriter, r *http.Request) { // why serviceHandler using pointer?
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var s model.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow(`
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

func (h *ServiceHandler) ServiceGet(w http.ResponseWriter, r *http.Request) { // r request for data and w writes or provide that data
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.DB.Query(`SELECT id, name, description FROM services;`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // why?

	services := make([]model.Service, 0)
	// The loop basically means:
	// "Give me the first row → scan it → put it in my slice.
	// Give me the next row → scan it → put it in my slice.
	// Keep going until there are no more rows."
	for rows.Next() {
		var s model.Service
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

func (h *ServiceHandler) ServiceId(w http.ResponseWriter, r *http.Request) {
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

	var s model.Service
	err = h.DB.QueryRow(`SELECT id, name, description FROM services WHERE id=$1;`, id).Scan(&s.ID, &s.Name, &s.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *ServiceHandler) ServicePut(w http.ResponseWriter, r *http.Request) {
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

	var s model.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(`
	UPDATE services
	SET name=$1, description=$2
	WHERE id=$3 RETURNING id, name, description`, s.Name, s.Description, id, // The SQL says what values it needs. Go supplies them.
	).Scan(&s.ID, &s.Name, &s.Description)
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

func (h *ServiceHandler) ServicePatch(w http.ResponseWriter, r *http.Request) {
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

	var s model.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(`
	UPDATE services
	SET name=COALESCE($1, name),
	description=COALESCE($2, description)
	WHERE id=$3
	RETURNING id, name, description`,
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

func (h *ServiceHandler) ServiceDelete(w http.ResponseWriter, r *http.Request) {
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
	err = h.DB.QueryRow(`DELETE FROM services WHERE id=$1 RETURNING id`, id).Scan(&deletedId)
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
