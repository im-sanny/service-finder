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

	pH := handler.ProviderHandler{
		DB: db,
	}

	mux.HandleFunc("POST /service", sH.ServicePost)
	mux.HandleFunc("GET /service", sH.ServiceGet)
	mux.HandleFunc("GET /service/{id}", sH.ServiceId)
	mux.HandleFunc("PUT /service/{id}", sH.ServicePut)
	mux.HandleFunc("PATCH /service/{id}", sH.ServicePatch)
	mux.HandleFunc("DELETE /service/{id}", sH.ServiceDelete)

	mux.HandleFunc("POST /provider", pH.ProviderPost)
	mux.HandleFunc("GET /provider", pH.ProviderGet)
	mux.HandleFunc("GET /provider/{id}", pH.ProviderID)
	mux.HandleFunc("PUT /provider/{id}", pH.ProviderPut)
	mux.HandleFunc("PATCH /provider/{id}", pH.ProviderPatch)
	mux.HandleFunc("DELETE /provider/{id}", pH.ProviderDelete)

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
