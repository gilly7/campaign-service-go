package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"campaign-service/internal/api"
	"campaign-service/internal/campaign"
	"campaign-service/internal/database"
	"campaign-service/internal/message"
)

func main() {
	db := database.Connect()
	defer db.Close()

	queue := message.NewRedisQueue()
	service := campaign.NewService(db, queue)
	handler := api.NewHandler(service)

	// Background worker
	go service.StartWorker(context.Background())

	// HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.Router(),
	}

	go func() {
		log.Println("Server starting on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("Service running. Press Ctrl+C to stop.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutdown complete.")
}
