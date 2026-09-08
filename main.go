package main

import (
	"log"
	"net/http"

	"github.com/im-sanny/service-finder/database"
	"github.com/im-sanny/service-finder/handler"
	"github.com/im-sanny/service-finder/repository"
)

func main() {
	conStr := "postgres://postgres:360420@localhost:5432/serfin?sslmode=disable"

	db, err := database.Connect(conStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	mux := http.NewServeMux()

	serviceRepo := repository.NewServiceRepository(db)
	sH := handler.NewServiceHandler(serviceRepo)

	providerRepo := repository.NewProviderRepository(db)
	pH := handler.NewProviderHandler(providerRepo)

	mux.HandleFunc("POST /services", sH.Create)
	mux.HandleFunc("GET /services", sH.GetAll)
	mux.HandleFunc("GET /services/{id}", sH.GetByID)
	mux.HandleFunc("PUT /services/{id}", sH.Update)
	mux.HandleFunc("PATCH /services/{id}", sH.Patch)
	mux.HandleFunc("DELETE /services/{id}", sH.Delete)

	mux.HandleFunc("POST /providers", pH.Create)
	mux.HandleFunc("GET /providers", pH.GetAll)
	mux.HandleFunc("GET /providers/{id}", pH.GetByID)
	mux.HandleFunc("PUT /providers/{id}", pH.Update)
	mux.HandleFunc("PATCH /providers/{id}", pH.Patch)
	mux.HandleFunc("DELETE /providers/{id}", pH.Delete)

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
