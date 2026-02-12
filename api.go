package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// API handles HTTP requests
type API struct {
	storage *Storage
	pool    *WorkerPool
}

// NewAPI creates a new API instance
func NewAPI(storage *Storage, pool *WorkerPool) *API {
	return &API{
		storage: storage,
		pool:    pool,
	}
}

// CreateJobHandler handles POST /api/jobs - Submit a new scraping job
func (api *API) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var req ScrapeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.URLs) == 0 {
		http.Error(w, "URLs list cannot be empty", http.StatusBadRequest)
		return
	}

	// Create a new job
	job := &ScrapeJob{
		ID:        uuid.New().String(),
		URLs:      req.URLs,
		Selector:  req.Selector,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save job
	api.storage.SaveJob(job)

	// Update status to running
	job.Status = StatusRunning
	job.UpdatedAt = time.Now()
	api.storage.SaveJob(job)

	// Submit tasks to worker pool
	for _, url := range job.URLs {
		api.pool.Submit(ScrapeTask{
			JobID:    job.ID,
			URL:      url,
			Selector: job.Selector,
		})
	}

	// Return job info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JobResponse{
		Job: *job,
	})
}

// GetJobHandler handles GET /api/jobs/{id} - Get job status and results
func (api *API) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]

	job, exists := api.storage.GetJob(jobID)

	if exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	results := api.storage.GetResults(jobID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JobResponse{
		Job:     *job,
		Results: results,
	})
}

// ListJobsHandler handles GET /api/jobs - List all jobs
func (api *API) ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs := api.storage.GetAllJobs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// GetResultsHandler handles GET /api/results - Get all results (with optional filtering)
func (api *API) GetResultsHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")

	var results []ScrapeResult

	if jobID != "" {
		results = api.storage.GetResults(jobID)
	}

	results = api.storage.GetAllResults()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"count":   len(results),
	})
}

// HealthHandler handles GET /api/health - Health check endpoint
func (api *API) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// SetupRoutes configures all API routes
func (api *API) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/health", api.HealthHandler).Methods("GET")
	router.HandleFunc("/api/jobs", api.CreateJobHandler).Methods("POST")
	router.HandleFunc("/api/jobs", api.ListJobsHandler).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", api.GetJobHandler).Methods("GET")
	router.HandleFunc("/api/results", api.GetResultsHandler).Methods("GET")

	return router
}
