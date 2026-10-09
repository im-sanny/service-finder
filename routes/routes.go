package routes

import (
	"net/http"

	"github.com/im-sanny/service-finder/handler"
	"github.com/im-sanny/service-finder/middleware"
)

// registerRoutes sets up all HTTP handlers and applies middleware
func RegisterRoutes(mux *http.ServeMux, sH *handler.ServiceHandler, pH *handler.ProviderHandler) {
	// --- Public Routes (Read-only) ---
	mux.HandleFunc("GET /services", sH.GetAll)
	mux.HandleFunc("GET /services/{id}", sH.GetByID)
	mux.HandleFunc("GET /providers", pH.GetAll)
	mux.HandleFunc("GET /providers/{id}", pH.GetByID)

	// --- Protected Routes (Write/Delete) ---
	// Helper to reduce repetition
	protected := func(method string, path string, handler http.HandlerFunc) {
		mux.Handle(method+" "+path, middleware.RequireAPIKey(handler))
	}

	// Services
	protected("POST", "/services", sH.Create)
	protected("PUT", "/services/{id}", sH.Update)
	protected("PATCH", "/services/{id}", sH.Patch)
	protected("DELETE", "/services/{id}", sH.Delete)
	protected("POST", "/services/batch", sH.CreateBatch)
	protected("DELETE", "/services/batch", sH.DeleteBatch)

	// Providers
	protected("POST", "/providers", pH.Create)
	protected("PUT", "/providers/{id}", pH.Update)
	protected("PATCH", "/providers/{id}", pH.Patch)
	protected("DELETE", "/providers/{id}", pH.Delete)
	protected("POST", "/providers/batch", pH.CreateBatch)
	protected("DELETE", "/providers/batch", pH.DeleteBatch)
}
