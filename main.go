package main

import (
	"log"

	"github.com/gauravpandey771/task-api/internal/domain"
	"github.com/gauravpandey771/task-api/internal/repository"
	httphandler "github.com/gauravpandey771/task-api/internal/transport/http"
)

func main() {
	// Initialize repository (in-memory)
	repo := repository.NewInMemoryTaskRepository()

	// Initialize service
	service := domain.NewTaskService(repo)

	// Initialize HTTP handler
	handler := httphandler.NewTaskHandler(service)

	// Create and start Fiber app
	app := httphandler.NewApp(handler)

	log.Println("Starting Task Management API on :8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
