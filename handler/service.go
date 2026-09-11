package handler

import (
	"net/http"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/service"
)

type ServiceHandler struct {
	svc service.Service
}

func NewServiceHandler(svc service.Service) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

// - *ServiceHandler: avoids copying the struct, shares the DB pool.
// - *http.Request: avoids copying large request data, allows body/context reading.
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s model.Service
	if !decodeJSON(w, r, &s) {
		return
	}

	if err := h.svc.Create(&s); err != nil {
		// Note: In a real app, you should log the actual 'err' here before returning a generic message
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, s)
}

func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	services, err := h.svc.GetAll()
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, services)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	s, err := h.svc.GetByID(id)
	if err != nil {
		respondError(w, err)
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

	if err := h.svc.Update(&s); err != nil {
		respondError(w, err)
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

	s, err := h.svc.Patch(id, update.Name, update.Description)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := getIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.svc.Delete(id); err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}
