package main

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// WorkerPool manages a pool of workers for concurrent scraping
type WorkerPool struct {
	workers     int
	jobQueue    chan ScrapeTask
	limiter     *rate.Limiter
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	scraper     *Scraper
	storage     *Storage
}

// ScrapeTask represents a task to be processed by the worker pool
type ScrapeTask struct {
	JobID    string
	URL      string
	Selector string
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, rateLimit int, scraper *Scraper, storage *Storage) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan ScrapeTask, workers*2), // Buffer for smooth operation
		limiter:  rate.NewLimiter(rate.Limit(rateLimit), rateLimit), // requests per second
		ctx:      ctx,
		cancel:   cancel,
		scraper:  scraper,
		storage:  storage,
	}
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the goroutine that processes tasks
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case task, ok := <-wp.jobQueue:
			if !ok {
				return
			}

			// Wait for rate limiter
			if err := wp.limiter.Wait(wp.ctx); err != nil {
				// Context cancelled
				return
			}

			// Process the scraping task
			result := wp.scraper.Scrape(task.URL, task.Selector)
			result.JobID = task.JobID
			result.ScrapedAt = time.Now()

			// Store the result
			wp.storage.SaveResult(result)

			// Check if all URLs for this job are complete
			wp.checkJobCompletion(task.JobID)
		}
	}
}

// Submit adds a task to the job queue
func (wp *WorkerPool) Submit(task ScrapeTask) {
	select {
	case wp.jobQueue <- task:
	case <-wp.ctx.Done():
	}
}

// checkJobCompletion checks if a job has completed all its URLs
func (wp *WorkerPool) checkJobCompletion(jobID string) {
	job, exists := wp.storage.GetJob(jobID)
	if !exists {
		return
	}

	results := wp.storage.GetResults(jobID)

	// If we have results for all URLs, mark job as completed
	if len(results) >= len(job.URLs) {
		job.Status = StatusCompleted
		job.UpdatedAt = time.Now()
		wp.storage.SaveJob(job)
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
}
