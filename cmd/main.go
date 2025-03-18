package main

import (
	"berezovskiyvalerii/csv-rest-app/internal/delivery/rest"
	"berezovskiyvalerii/csv-rest-app/internal/repository"
	"berezovskiyvalerii/csv-rest-app/internal/server"
	"berezovskiyvalerii/csv-rest-app/internal/service"
	"berezovskiyvalerii/csv-rest-app/pkg/db"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	// "github.com/spf13/viper"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Errorf("error loading env file: %s", err.Error())
		return
	}

	log.Info("env loaded successfully")

	//Database initialization
	db, err := db.NewPostgresConnection(db.ConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("could not initialize database: %s", err.Error())
	}
	defer db.Close()
	log.Info("DB initialized")

	repo := repository.NewProducts(db)
	service := service.NewProducts(repo)
	handler := rest.NewHandler(service)

	srv := new(server.Server)
	go func() {
		//TODO: use lib for 'port'
		if err := srv.Run("8000", handler.InitRoutes()); err != nil {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	log.Info("Server started")

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("Error closing database connection: %s", err.Error())
	}

	log.Info("Server exited gracefully")
}
