package handler

import (
	"database/sql"
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
	var s model.Service
	if !decodeJSON(w, r, &s) {
		return
	}

	if err := h.repo.Create(&s); err != nil {
		// Note: In a real app, you should log the actual 'err' here before returning a generic message
		http.Error(w, "Failed to create service", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	services, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, services)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, s)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var s model.Service
	if !decodeJSON(w, r, &s) {
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

	writeJSON(w, http.StatusOK, s)
}

func (h *ServiceHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var update model.Service
	if !decodeJSON(w, r, &update) {
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

	writeJSON(w, http.StatusOK, s)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusNoContent, nil)
}
