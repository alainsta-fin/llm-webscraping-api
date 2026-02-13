package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Scraper handles web scraping operations
type Scraper struct {
	client  *http.Client
	timeout time.Duration
}

// NewScraper creates a new Scraper instance
func NewScraper(timeout time.Duration) *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// Scrape fetches and parses a URL
func (s *Scraper) Scrape(url, selector string) ScrapeResult {
	result := ScrapeResult{
		URL:     url,
		Success: false,
	}

	// Make HTTP request
	resp, err := s.client.Get(url)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to fetch URL: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("Non-200 status code: %d", resp.StatusCode)
		return result
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read response: %v", err)
		return result
	}

	htmlContent := string(body)

	// Extract title
	result.Title = s.extractTitle(htmlContent)

	// Extract content based on selector or full text
	if selector != "" {
		result.Content = s.extractBySelector(htmlContent, selector)
	} else {
		result.Content = s.extractText(htmlContent)
	}

	// Extract links
	result.Links = s.extractLinks(htmlContent, url)

	result.Success = true
	return result
}

// extractTitle extracts the page title from HTML
func (s *Scraper) extractTitle(html string) string {
	re := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// extractText extracts visible text from HTML (simplified)
func (s *Scraper) extractText(html string) string {
	// Remove script and style tags
	re := regexp.MustCompile(`(?s)<(script|style)[^>]*>.*?</\1>`)
	html = re.ReplaceAllString(html, "")

	// Remove HTML tags
	re = regexp.MustCompile(`<[^>]+>`)
	text := re.ReplaceAllString(html, " ")

	// Clean up whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	text = strings.TrimSpace(text)

	// Limit content length
	if len(text) > 5000 {
		text = text[:5000] + "..."
	}

	return text
}

// extractBySelector extracts content matching a simple selector (simplified CSS selector support)
func (s *Scraper) extractBySelector(html, selector string) string {
	// Basic support for class and id selectors
	var pattern string

	if strings.HasPrefix(selector, ".") {
		// Class selector
		className := strings.TrimPrefix(selector, ".")
		pattern = fmt.Sprintf(`(?s)<[^>]*class="[^"]*%s[^"]*"[^>]*>(.*?)</[^>]+>`, regexp.QuoteMeta(className))
	} else if strings.HasPrefix(selector, "#") {
		// ID selector
		id := strings.TrimPrefix(selector, "#")
		pattern = fmt.Sprintf(`(?s)<[^>]*id="%s"[^>]*>(.*?)</[^>]+>`, regexp.QuoteMeta(id))
	} else {
		// Tag selector
		pattern = fmt.Sprintf(`(?s)<%s[^>]*>(.*?)</%s>`, regexp.QuoteMeta(selector), regexp.QuoteMeta(selector))
	}

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(html, -1)

	var results []string
	for _, match := range matches {
		if len(match) > 1 {
			text := s.extractText(match[1])
			if text != "" {
				results = append(results, text)
			}
		}
	}

	content := strings.Join(results, "\n\n")

	// Limit content length
	if len(content) > 5000 {
		content = content[:5000] + "..."
	}

	return content
}

// extractLinks extracts all links from HTML
func (s *Scraper) extractLinks(html, baseURL string) []string {
	re := regexp.MustCompile(`<a[^>]+href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(html, -1)

	links := make([]string, 0)
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]

			// Skip anchors, javascript, and mailto links
			if strings.HasPrefix(link, "#") ||
				strings.HasPrefix(link, "javascript:") ||
				strings.HasPrefix(link, "mailto:") {
				continue
			}

			// Make relative URLs absolute (simplified)
			if strings.HasPrefix(link, "/") {
				// Extract base domain from baseURL
				re := regexp.MustCompile(`^(https?://[^/]+)`)
				if domainMatch := re.FindStringSubmatch(baseURL); len(domainMatch) > 1 {
					link = domainMatch[1] + link
				}
			}

			// Avoid duplicates
			// BUG HERE...remove !seen[link]
			if !seen[link] && (strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")) {
				seen[link] = true
				links = append(links, link)

				// Bug Here...remove number of links
				// Limit number of links
				if len(links) >= 50 {
					break
				}
			}
		}
	}

	return links
}
