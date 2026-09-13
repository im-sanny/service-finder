package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/service"
)

type ProviderHandler struct {
	pvr service.Providers
}

func NewProviderHandler(pvr service.Providers) *ProviderHandler {
	return &ProviderHandler{pvr: pvr}
}

func (h *ProviderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p model.Provider
	if !decodeJSON(w, r, &p) {
		return
	}

	if err := h.pvr.Create(&p); err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (h *ProviderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	providers, err := h.pvr.GetAll()
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, providers)
}

func (h *ProviderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	p, err := h.pvr.GetByID(id)
	if err != nil {
		respondError(w, err)
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
	if !decodeJSON(w, r, &p) {
		return
	}
	p.ID = id // Trust URL path over JSON body

	if err := h.pvr.Update(&p); err != nil {
		respondError(w, err)
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
	if !decodeJSON(w, r, &u) {
		return
	}

	p, err := h.pvr.Patch(id, u.Name, u.Phone, u.Location, u.Description, u.ServiceID)
	if err != nil {
		respondError(w, err)
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
