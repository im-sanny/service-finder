package middleware

import (
	"net/http"
	"os"
)

// RequireAPIKey is a middleware that checks for a valid API key in the header.
func RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the key from the environment variable
		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		// 2. Check the "X-API-Key" header in the request
		apiKey := r.Header.Get("X-API-Key")

		// 3. Compare
		if apiKey != expectedKey {
			http.Error(w, "Unauthorized: Invalid or missing API key", http.StatusUnauthorized)
			return
		}

		// 4. If valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
