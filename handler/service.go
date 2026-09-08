package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/repository"
)

type ServiceHandler struct {
	repo repository.ServiceRepository
}

func NewServiceHandler(repo repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{repo: repo}
}

// - *ServiceHandler: avoids copying the struct, shares the DB pool.
// - *http.Request: avoids copying large request data, allows body/context reading.
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var s model.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&s); err != nil {
		// Note: In a real app, you should log the actual 'err' here before returning a generic message
		http.Error(w, "Failed to create service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	s, err := h.repo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var s model.Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	s.ID = id

	if err := h.repo.Update(&s); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *ServiceHandler) Patch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var update model.Service
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	s, err := h.repo.Patch(id, update.Name, update.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to patch service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
