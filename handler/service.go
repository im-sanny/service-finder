package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/service"
)

type ServiceHandler struct {
	svc service.Service
}

func NewServiceHandler(svc service.Service) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

func (h *ServiceHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	var svc []*model.Service
	if !decodeJSON(w, r, &svc) {
		return
	}

	create, err := h.svc.CreateBatch(svc)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, create)
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
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	services, total, err := h.svc.GetAll(page, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	response := map[string]any{
		"data":  services,
		"total": total,
		"page":  page,
		"limit": limit,
	}
	
	writeJSON(w, http.StatusOK, response)
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

func (h *ServiceHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	// Parse comma-separated IDs from query string
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "ids query parameter is required", http.StatusBadRequest)
		return
	}

	// Split and parse IDs
	idString := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(idString))
	for i, idStr := range idString {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			http.Error(w, "invalid id format at position"+strconv.Itoa(i), http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	// call service layer
	deleted, err := h.svc.DeleteBatch(ids)
	if err != nil {
		respondError(w, err)
		return
	}

	// Return count of deleted row
	response := map[string]int64{"deleted": deleted}
	writeJSON(w, http.StatusOK, response)
}
