package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
func (h *ServiceHandler) ServicePost(w http.ResponseWriter, r *http.Request) {
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

	services, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
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

	s, err := h.repo.GetById(int64(id))
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

	s.ID = int64(id)

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

	var update model.Service
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	s, err := h.repo.Patch(int64(id), *update.Name, *update.Description)
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

	if err := h.repo.Delete(int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
