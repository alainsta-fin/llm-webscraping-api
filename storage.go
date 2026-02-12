package main

import (
	"sync"
)

// Storage provides thread-safe storage for jobs and results
type Storage struct {
	jobs    map[string]*ScrapeJob
	results map[string][]ScrapeResult // jobID -> results
	mu      sync.RWMutex
}

// NewStorage creates a new Storage instance
func NewStorage() *Storage {
	return &Storage{
		jobs:    make(map[string]*ScrapeJob),
		results: make(map[string][]ScrapeResult),
	}
}

// SaveJob stores or updates a job
func (s *Storage) SaveJob(job *ScrapeJob) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

// GetJob retrieves a job by ID
func (s *Storage) GetJob(jobID string) (*ScrapeJob, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, exists := s.jobs[jobID]
	return job, exists
}

// GetAllJobs retrieves all jobs
func (s *Storage) GetAllJobs() []*ScrapeJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	jobs := make([]*ScrapeJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// SaveResult stores a scrape result
func (s *Storage) SaveResult(result ScrapeResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.JobID] = append(s.results[result.JobID], result)
}

// GetResults retrieves all results for a job
func (s *Storage) GetResults(jobID string) []ScrapeResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results, exists := s.results[jobID]
	if !exists {
		return []ScrapeResult{}
	}
	return results
}

// GetAllResults retrieves all results across all jobs
func (s *Storage) GetAllResults() []ScrapeResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var allResults []ScrapeResult
	for _, results := range s.results {
		allResults = append(allResults, results...)
	}
	return allResults
}
