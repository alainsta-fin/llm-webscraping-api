package main

import (
	"time"
)

// JobStatus represents the status of a scraping job
type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

// ScrapeRequest represents a request to scrape a URL
type ScrapeRequest struct {
	URLs     []string `json:"urls"`
	Selector string   `json:"selector,omitempty"` // CSS selector for targeted scraping
}

// ScrapeJob represents a scraping job
type ScrapeJob struct {
	ID        string    `json:"id"`
	URLs      []string  `json:"urls"`
	Selector  string    `json:"selector,omitempty"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScrapeResult represents the result of scraping a single URL
type ScrapeResult struct {
	JobID      string    `json:"job_id"`
	URL        string    `json:"url"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	Links      []string  `json:"links,omitempty"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	Error      string    `json:"error,omitempty"`
	ScrapedAt  time.Time `json:"scraped_at"`
}

// JobResponse represents the response when creating/querying a job
type JobResponse struct {
	Job     ScrapeJob      `json:"job"`
	Results []ScrapeResult `json:"results,omitempty"`
}
