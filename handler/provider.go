package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/im-sanny/service-finder/model"
)

type ProviderHandler struct {
	DB *sql.DB
}

func (h *ProviderHandler) ProviderPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var p model.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	err := h.DB.QueryRow(`INSERT INTO providers (name, phone, location, description, service_id) VALUES($1, $2, $3, $4, $5) RETURNING id`, p.Name, p.Phone, p.Location, p.Description, p.ServiceID).Scan(&p.ID)

	if err != nil {
		http.Error(w, "Failed to insert service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProviderHandler) ProviderGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.DB.Query(`SELECT id, name, phone, location, description, service_id FROM providers`)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var provider []model.Provider
	for rows.Next() {
		var p model.Provider
		if err := rows.Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		provider = append(provider, p)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Row iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)

}

func (h *ProviderHandler) ProviderID(w http.ResponseWriter, r *http.Request) {
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

	var p model.Provider
	err = h.DB.QueryRow(`SELECT id, name, phone, location, description, service_id FROM providers WHERE id=$1`,
		id).Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProviderHandler) ProviderPut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p model.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(`
	UPDATE providers SET name=$1, phone=$2, location=$3, description=$4, service_id=$5 WHERE id=$6
	RETURNING id, name, phone, location, description, service_id`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID, id).Scan(&p.ID, &p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProviderHandler) ProviderPatch(w http.ResponseWriter, r *http.Request) {
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

	var p model.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.DB.QueryRow(`
	UPDATE providers SET
	name=COALESCE($1, name),
	phone=COALESCE($2, phone),
	location=COALESCE($3, location),
	description=COALESCE($4, description),
	service_id=COALESCE($5,	service_id)
	WHERE id=$6 RETURNING id, name, phone, location, description, service_id`,
		p.Name, p.Phone, p.Location, p.Description, p.ServiceID, id).Scan(&p.ID,
		&p.Name, &p.Phone, &p.Location, &p.Description, &p.ServiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProviderHandler) ProviderDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Use h.DB.Exec for operations that do not return rows
	res, err := h.DB.Exec(`DELETE FROM providers WHERE id=$1`, id)
	if err != nil {
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	// Check how many rows were actually deleted
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to verify deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204 No Content is the standard REST response
}
