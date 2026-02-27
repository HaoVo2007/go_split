package main

import (
	"go-split/internal/app"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", container.Config.Server.Host, container.Config.Server.Port)
	log.Printf("Server starting on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := container.Router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := container.MongoDB.Client().Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	}

	log.Println("Server exited gracefully")
}
