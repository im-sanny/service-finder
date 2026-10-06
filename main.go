package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/im-sanny/service-finder/database"
	"github.com/im-sanny/service-finder/handler"
	"github.com/im-sanny/service-finder/middleware"
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

	// public routers
	mux.HandleFunc("GET /services", sH.GetAll)
	mux.HandleFunc("GET /services/{id}", sH.GetByID)
	mux.HandleFunc("GET /providers", pH.GetAll)
	mux.HandleFunc("GET /providers/{id}", pH.GetByID)

	// protected routers
	mux.Handle("POST /services/batch", middleware.RequireAPIKey(http.HandlerFunc(sH.CreateBatch)))
	mux.Handle("POST /services", middleware.RequireAPIKey(http.HandlerFunc(sH.Create)))
	mux.Handle("PUT /services/{id}", middleware.RequireAPIKey(http.HandlerFunc(sH.Update)))
	mux.Handle("PATCH /services/{id}", middleware.RequireAPIKey(http.HandlerFunc(sH.Patch)))
	mux.Handle("DELETE /services/{id}", middleware.RequireAPIKey(http.HandlerFunc(sH.Delete)))
	mux.Handle("DELETE /services/batch", middleware.RequireAPIKey(http.HandlerFunc(sH.DeleteBatch)))

	mux.Handle("POST /providers/batch", middleware.RequireAPIKey(http.HandlerFunc(pH.CreateBatch)))
	mux.Handle("POST /providers", middleware.RequireAPIKey(http.HandlerFunc(pH.Create)))
	mux.Handle("PUT /providers/{id}", middleware.RequireAPIKey(http.HandlerFunc(pH.Update)))
	mux.Handle("PATCH /providers/{id}", middleware.RequireAPIKey(http.HandlerFunc(pH.Patch)))
	mux.Handle("DELETE /providers/{id}", middleware.RequireAPIKey(http.HandlerFunc(pH.Delete)))
	mux.Handle("DELETE /providers/batch", middleware.RequireAPIKey(http.HandlerFunc(pH.DeleteBatch)))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server running on port :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")

	// give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Printf("Server exited properly")
}
