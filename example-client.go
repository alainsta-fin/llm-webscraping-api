package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Example client demonstrating how to use the Web Scraper API

const apiURL = "http://localhost:8080"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SubmitJob(urls []string, selector string) (string, error) {
	reqBody := map[string]interface{}{
		"urls":     urls,
		"selector": selector,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/jobs",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	job := result["job"].(map[string]interface{})
	jobID := job["id"].(string)

	return jobID, nil
}

func (c *Client) GetJob(jobID string) (map[string]interface{}, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/jobs/" + jobID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetResults(jobID string) ([]interface{}, error) {
	url := c.baseURL + "/api/results"
	if jobID != "" {
		url += "?job_id=" + jobID
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	results := result["results"].([]interface{})
	return results, nil
}

func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/api/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Health check: %s\n", string(body))

	return nil
}

// Example usage (commented out - this is meant to show how to use the client)
/*
func main() {
	client := NewClient("http://localhost:8080")

	// Health check
	fmt.Println("=== Health Check ===")
	if err := client.HealthCheck(); err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		return
	}
	fmt.Println()

	// Submit a job
	fmt.Println("=== Submitting Job ===")
	urls := []string{
		"https://example.com",
		"https://golang.org",
	}

	jobID, err := client.SubmitJob(urls, "")
	if err != nil {
		fmt.Printf("Failed to submit job: %v\n", err)
		return
	}
	fmt.Printf("Job submitted successfully! Job ID: %s\n", jobID)
	fmt.Println()

	// Wait for job to complete
	fmt.Println("=== Waiting for Job to Complete ===")
	time.Sleep(5 * time.Second)

	// Get job status
	fmt.Println("=== Getting Job Status ===")
	job, err := client.GetJob(jobID)
	if err != nil {
		fmt.Printf("Failed to get job: %v\n", err)
		return
	}

	jobData := job["job"].(map[string]interface{})
	fmt.Printf("Job Status: %s\n", jobData["status"])
	fmt.Println()

	// Get results
	fmt.Println("=== Getting Results ===")
	results, err := client.GetResults(jobID)
	if err != nil {
		fmt.Printf("Failed to get results: %v\n", err)
		return
	}

	for i, result := range results {
		r := result.(map[string]interface{})
		fmt.Printf("\nResult %d:\n", i+1)
		fmt.Printf("  URL: %s\n", r["url"])
		fmt.Printf("  Title: %s\n", r["title"])
		fmt.Printf("  Success: %v\n", r["success"])
		fmt.Printf("  Status Code: %.0f\n", r["status_code"])

		if links, ok := r["links"].([]interface{}); ok {
			fmt.Printf("  Links found: %d\n", len(links))
		}
	}
}
*/
