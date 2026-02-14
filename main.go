package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// Server configuration
	DefaultPort        = "8080"
	DefaultWorkers     = 10
	DefaultRateLimit   = 5  // requests per second
	DefaultHTTPTimeout = 30 // seconds
)

func main() {
	// Get configuration from environment or use defaults
	port := getEnv("PORT", DefaultPort)
	workers := DefaultWorkers
	rateLimit := DefaultRateLimit
	httpTimeout := time.Duration(DefaultHTTPTimeout) * time.Second

	// Initialize components
	storage := NewStorage()
	scraper := NewScraper(httpTimeout)
	pool := NewWorkerPool(workers, rateLimit, scraper, storage)
	api := NewAPI(storage, pool)

	// Start worker pool
	pool.Start()
	log.Printf("Worker pool started with %d workers (rate limit: %d req/s)", workers, rateLimit)

	// Setup HTTP server
	router := api.SetupRoutes()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		log.Printf("API endpoints:")
		log.Printf("  POST   /api/jobs       - Submit scraping job")
		log.Printf("  GET    /api/jobs       - List all jobs")
		log.Printf("  GET    /api/jobs/{id}  - Get job status and results")
		log.Printf("  GET    /api/results    - Get all results (optional: ?job_id=xxx)")
		log.Printf("  GET    /api/health     - Health check")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
