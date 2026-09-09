package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
