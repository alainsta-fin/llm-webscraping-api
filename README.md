# Concurrent Web Scraper API

A high-performance concurrent web scraper API built with Go, featuring goroutine pools, rate limiting, and RESTful endpoints.

## Features

- **REST API** for submitting and managing scraping jobs
- **Goroutine pool** with configurable workers for parallel scraping
- **Rate limiting** to avoid overwhelming target servers
- **Structured storage** for jobs and results
- **Query endpoints** for retrieving scraped data
- **Graceful shutdown** handling

## Architecture

- `main.go` - Server initialization and configuration
- `api.go` - REST API handlers and routes
- `models.go` - Data structures for jobs and results
- `storage.go` - Thread-safe in-memory storage
- `pool.go` - Worker pool with rate limiting
- `scraper.go` - Web scraping logic

## Installation

```bash
# Clone or navigate to the project directory
cd web-scraper-api

# Download dependencies
go mod download

# Build the application
go build -o scraper-api

# Run the server
./scraper-api
```

Or run directly:
```bash
go run .
```

## Configuration

Environment variables:
- `PORT` - Server port (default: 8080)

Hardcoded defaults (can be modified in main.go):
- Workers: 10 concurrent goroutines
- Rate limit: 5 requests per second
- HTTP timeout: 30 seconds

## API Endpoints

### 1. Submit Scraping Job

**POST** `/api/jobs`

Submit URLs to be scraped.

Request body:
```json
{
  "urls": [
    "https://example.com",
    "https://golang.org"
  ],
  "selector": ".main-content"
}
```

Response:
```json
{
  "job": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "urls": ["https://example.com", "https://golang.org"],
    "selector": ".main-content",
    "status": "running",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Get Job Status and Results

**GET** `/api/jobs/{id}`

Retrieve job information and all results.

Response:
```json
{
  "job": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "completed",
    "urls": ["https://example.com"],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:05Z"
  },
  "results": [
    {
      "job_id": "123e4567-e89b-12d3-a456-426614174000",
      "url": "https://example.com",
      "title": "Example Domain",
      "content": "Example text content...",
      "links": ["https://example.com/page1"],
      "status_code": 200,
      "success": true,
      "scraped_at": "2024-01-15T10:30:05Z"
    }
  ]
}
```

### 3. List All Jobs

**GET** `/api/jobs`

List all scraping jobs.

Response:
```json
{
  "jobs": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "status": "completed",
      "urls": ["https://example.com"],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:05Z"
    }
  ],
  "count": 1
}
```

### 4. Query Results

**GET** `/api/results`

Get all results, optionally filtered by job ID.

Query parameters:
- `job_id` (optional) - Filter results by job ID

Response:
```json
{
  "results": [
    {
      "job_id": "123e4567-e89b-12d3-a456-426614174000",
      "url": "https://example.com",
      "title": "Example Domain",
      "content": "...",
      "success": true,
      "status_code": 200
    }
  ],
  "count": 1
}
```

### 5. Health Check

**GET** `/api/health`

Check if the API is running.

Response:
```json
{
  "status": "ok",
  "time": "2024-01-15T10:30:00Z"
}
```

## Usage Examples

### Using cURL

```bash
# Submit a scraping job
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://example.com", "https://golang.org"],
    "selector": "body"
  }'

# Get job status (replace {job_id} with actual ID from previous response)
curl http://localhost:8080/api/jobs/{job_id}

# List all jobs
curl http://localhost:8080/api/jobs

# Get all results
curl http://localhost:8080/api/results

# Get results for specific job
curl "http://localhost:8080/api/results?job_id={job_id}"

# Health check
curl http://localhost:8080/api/health
```

### CSS Selectors

The `selector` field supports basic CSS selectors:

- Tag selector: `"div"`, `"p"`, `"article"`
- Class selector: `".main-content"`, `".article-body"`
- ID selector: `"#content"`, `"#main"`
- Leave empty for full page text extraction

## Job Status

- `pending` - Job created but not yet started
- `running` - Scraping in progress
- `completed` - All URLs scraped successfully
- `failed` - Job encountered errors

## Features Details

### Concurrent Processing
- Configurable number of worker goroutines
- Tasks distributed across workers via buffered channel
- Automatic load balancing

### Rate Limiting
- Token bucket algorithm (golang.org/x/time/rate)
- Prevents overwhelming target servers
- Configurable requests per second

### Data Extraction
- Page title extraction
- Full text content (with script/style removal)
- CSS selector-based content extraction
- Link extraction (up to 50 links per page)
- Automatic relative-to-absolute URL conversion

### Storage
- Thread-safe in-memory storage
- Separate storage for jobs and results
- Efficient concurrent access with RWMutex

## Limitations

- In-memory storage (data lost on restart)
- Simplified CSS selector support
- Content limited to 5000 characters
- Maximum 50 links extracted per page
- No authentication for scraped sites
- No JavaScript rendering

## Future Enhancements

- Persistent storage (PostgreSQL, MongoDB)
- Advanced CSS selector support (using goquery)
- JavaScript rendering (using chromedp)
- Authentication support
- Webhook notifications
- Job scheduling and retries
- Export to JSON/CSV
- Pagination for large result sets

## License

MIT
