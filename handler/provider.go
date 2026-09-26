package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/im-sanny/service-finder/model"
	"github.com/im-sanny/service-finder/service"
)

type ProviderHandler struct {
	pvr service.Providers
}

func NewProviderHandler(pvr service.Providers) *ProviderHandler {
	return &ProviderHandler{pvr: pvr}
}

func (h *ProviderHandler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	var providers []*model.Provider
	if !decodeJSON(w, r, &providers) {
		return
	}

	created, err := h.pvr.CreateBatch(providers)
	if err != nil {
		respondError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
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
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	filters := make(map[string]string)
	if loc := r.URL.Query().Get("location"); loc != "" {
		filters["location"] = loc
	}

	if sid := r.URL.Query().Get("service_id"); sid != "" {
		filters["service_id"] = sid
	}

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = 1
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = 10
	}

	providers, total, err := h.pvr.GetAll(page, limit, filters)
	if err != nil {
		respondError(w, err)
		return
	}

	response := map[string]any{
		"data":  providers,
		"total": total,
		"page":  page,
		"limit": limit,
	}

	writeJSON(w, http.StatusOK, response)
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
	if err := h.pvr.Delete(id); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProviderHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	// Parse comma-separated IDs from query string
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "ids query parameter is required", http.StatusBadRequest)
		return
	}

	// Split and parse IDs
	idStrings := strings.Split(idsParam, ",")
	ids := make([]int64, 0, len(idStrings))
	for i, idStr := range idStrings {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			http.Error(w, "invalid id format at position "+strconv.Itoa(i), http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	// Call service layer
	deleted, err := h.pvr.DeleteBatch(ids)
	if err != nil {
		respondError(w, err)
		return
	}

	// Return count of deleted rows
	response := map[string]int64{"deleted": deleted}
	writeJSON(w, http.StatusOK, response)
}
