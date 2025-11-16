package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	pg "pay-service/internal/adapters/postgres"

	ports "pay-service/internal/ports"
	pay_service "pay-service/internal/services"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := loadEnv(); err != nil {
		log.Fatalf("env: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	fmt.Printf("db url %v", dbURL)

	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	if err := pg.RunMigrations(dbURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations: OK")

	pool := pg.MustPool(ctx)
	defer pool.Close()
	log.Println("postgres: connected OK")

	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	service := pay_service.NewService()   
	handler := ports.NewHandler(service)

	api := server.Group("/api/v1")
	ports.Routers(api, handler)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           server,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()
	log.Println("http: listening on :8080")

	<-ctx.Done()
	log.Println("shutdown: start")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx) 

	log.Println("shutdown: done")
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		var pathErr *os.PathError
		switch {
		case errors.As(err, &pathErr):
			return nil
		default:
			return err
		}
	}
	return nil
}
