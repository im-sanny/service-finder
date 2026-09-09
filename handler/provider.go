package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/repository"
)

type ProviderHandler struct {
	repo repository.ProviderRepository
}

func NewProviderHandler(repo repository.ProviderRepository) *ProviderHandler {
	return &ProviderHandler{repo: repo}
}

func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p model.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&p); err != nil {
		http.Error(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	providers, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	p, err := h.repo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *ProviderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var p model.Provider
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	p.ID = id // Trust URL path over JSON body

	if err := h.repo.Update(&p); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update provider", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *ProviderHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	var u model.Provider
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	p, err := h.repo.Patch(id, u.Name, u.Phone, u.Location, u.Description, u.ServiceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to patch provider", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *ProviderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Provider not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
