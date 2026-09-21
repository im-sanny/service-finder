package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/im-sanny/service-finder/database"
	"github.com/im-sanny/service-finder/handler"
	"github.com/im-sanny/service-finder/repository"
	"github.com/im-sanny/service-finder/service"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

// 1. THIS LINE IS REQUIRED!
// It tells Go to bundle all .sql files from the migrations folder into embedMigration.
//
//go:embed migrations/*.sql
var embedMigration embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	conStr := os.Getenv("DB_URL")
	if conStr == "" {
		log.Fatal("DB_URL environment variable is not set. Please check your .env file.")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.Connect(conStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	goose.SetBaseFS(embedMigration)
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Printf("Database migrations applied successfully!")

	mux := http.NewServeMux()

	serviceRepo := repository.NewServiceRepository(db)
	svr := service.NewService(serviceRepo)
	sH := handler.NewServiceHandler(svr)

	providerRepo := repository.NewProviderRepository(db)
	pvr := service.NewProvider(providerRepo, serviceRepo)
	pH := handler.NewProviderHandler(pvr)

	mux.HandleFunc("POST /services/batch", sH.CreateBatch)
	mux.HandleFunc("POST /services", sH.Create)
	mux.HandleFunc("GET /services", sH.GetAll)
	mux.HandleFunc("GET /services/{id}", sH.GetByID)
	mux.HandleFunc("PUT /services/{id}", sH.Update)
	mux.HandleFunc("PATCH /services/{id}", sH.Patch)
	mux.HandleFunc("DELETE /services/{id}", sH.Delete)
	mux.HandleFunc("DELETE /services/batch/{id}", sH.DeleteBatch)

	mux.HandleFunc("POST /providers/batch", pH.CreateBatch)
	mux.HandleFunc("POST /providers", pH.Create)
	mux.HandleFunc("GET /providers", pH.GetAll)
	mux.HandleFunc("GET /providers/{id}", pH.GetByID)
	mux.HandleFunc("PUT /providers/{id}", pH.Update)
	mux.HandleFunc("PATCH /providers/{id}", pH.Patch)
	mux.HandleFunc("DELETE /providers/{id}", pH.Delete)
	mux.HandleFunc("DELETE /providers/batch/{id}", pH.DeleteBatch)

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
