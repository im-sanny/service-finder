package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/im-sanny/service-finder/internal/config"
	"github.com/im-sanny/service-finder/internal/database"
	"github.com/im-sanny/service-finder/internal/handler"
	"github.com/im-sanny/service-finder/internal/repository"
	"github.com/im-sanny/service-finder/internal/routes"
	"github.com/im-sanny/service-finder/internal/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	mux := http.NewServeMux()

	serviceRepo := repository.NewServiceRepository(db)
	svr := service.NewService(serviceRepo)
	sH := handler.NewServiceHandler(svr)

	providerRepo := repository.NewProviderRepository(db)
	pvr := service.NewProvider(providerRepo, serviceRepo)
	pH := handler.NewProviderHandler(pvr)

	userRepo := repository.NewUserRepository(db)
	uvr := service.NewAuthService(userRepo, cfg.Auth.JWTSecret)
	uH := handler.NewUserHandler(uvr)

	routes.RegisterRoutes(mux, sH, pH, uH)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Server running on port :%s", cfg.Server.Port)
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
