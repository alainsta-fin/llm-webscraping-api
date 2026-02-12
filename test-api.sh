#!/bin/bash

# Test script for Web Scraper API

API_URL="http://localhost:8080"

echo "========================================="
echo "Web Scraper API Test Suite"
echo "========================================="
echo ""

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. Health check
echo -e "${BLUE}1. Testing health endpoint...${NC}"
curl -s "${API_URL}/api/health" | jq '.'
echo ""
echo ""

# 2. Submit a scraping job
echo -e "${BLUE}2. Submitting scraping job...${NC}"
RESPONSE=$(curl -s -X POST "${API_URL}/api/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com",
      "https://golang.org"
    ],
    "selector": ""
  }')

echo "$RESPONSE" | jq '.'
JOB_ID=$(echo "$RESPONSE" | jq -r '.job.id')
echo -e "${GREEN}Job ID: $JOB_ID${NC}"
echo ""
echo ""

# 3. Wait a moment for scraping to complete
echo -e "${BLUE}3. Waiting for scraping to complete (5 seconds)...${NC}"
sleep 5
echo ""

# 4. Get job status
echo -e "${BLUE}4. Checking job status...${NC}"
curl -s "${API_URL}/api/jobs/${JOB_ID}" | jq '.'
echo ""
echo ""

# 5. List all jobs
echo -e "${BLUE}5. Listing all jobs...${NC}"
curl -s "${API_URL}/api/jobs" | jq '.'
echo ""
echo ""

# 6. Get all results
echo -e "${BLUE}6. Getting all results...${NC}"
curl -s "${API_URL}/api/results" | jq '.results[] | {url, title, success, status_code}'
echo ""
echo ""

# 7. Get results for specific job
echo -e "${BLUE}7. Getting results for job ${JOB_ID}...${NC}"
curl -s "${API_URL}/api/results?job_id=${JOB_ID}" | jq '.results[] | {url, title, success}'
echo ""
echo ""

# 8. Submit another job with selector
echo -e "${BLUE}8. Submitting job with CSS selector...${NC}"
RESPONSE2=$(curl -s -X POST "${API_URL}/api/jobs" \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com"
    ],
    "selector": "div"
  }')

echo "$RESPONSE2" | jq '.'
JOB_ID2=$(echo "$RESPONSE2" | jq -r '.job.id')
echo -e "${GREEN}Job ID: $JOB_ID2${NC}"
echo ""
echo ""

# 9. Wait and check second job
echo -e "${BLUE}9. Waiting for second job (3 seconds)...${NC}"
sleep 3
curl -s "${API_URL}/api/jobs/${JOB_ID2}" | jq '.results[] | {url, content: .content[:100]}'
echo ""
echo ""

echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}Test suite completed!${NC}"
echo -e "${GREEN}=========================================${NC}"
