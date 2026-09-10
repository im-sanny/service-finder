package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/im-sanny/service-finder/service"
)

// getIDFromPath extracts and validates the int64 ID from the URL path.
// This removes duplicated parsing logic across the handlers below.
func getIDFromPath(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	// ParseInt is safer than Atoi for int64 and prevents 32-bit overflow
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, data any) bool {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		http.Error(w, "Not found", http.StatusNotFound)
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		log.Printf("ERROR :%v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
