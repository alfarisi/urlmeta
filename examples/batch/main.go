package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alfarisi/urlmeta"
)

// Result holds the extraction result for a URL
type Result struct {
	URL      string
	Metadata *urlmeta.Metadata
	Error    error
	Duration time.Duration
}

func main() {
	fmt.Println("URLMeta - Batch Processing Example")
	fmt.Println("===================================\n")

	// List of URLs to process
	urls := []string{
		"https://github.com/golang/go",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://dev.to",
		"https://stackoverflow.com",
		"https://medium.com",
		"https://www.reddit.com",
		"https://twitter.com",
		"https://www.producthunt.com",
	}

	fmt.Printf("Processing %d URLs concurrently...\n\n", len(urls))

	// Sequential processing
	fmt.Println("=== Sequential Processing ===")
	startSeq := time.Now()
	resultsSeq := processSequential(urls)
	durationSeq := time.Since(startSeq)
	displayResults(resultsSeq)
	fmt.Printf("\nTotal time (sequential): %v\n", durationSeq)

	// Concurrent processing
	fmt.Println("\n=== Concurrent Processing ===")
	startConc := time.Now()
	resultsConc := processConcurrent(urls, 4) // 4 workers
	durationConc := time.Since(startConc)
	displayResults(resultsConc)
	fmt.Printf("\nTotal time (concurrent): %v\n", durationConc)

	// Performance comparison
	fmt.Println("\n=== Performance Comparison ===")
	fmt.Printf("Sequential: %v\n", durationSeq)
	fmt.Printf("Concurrent: %v\n", durationConc)
	speedup := float64(durationSeq) / float64(durationConc)
	fmt.Printf("Speedup: %.2fx faster\n", speedup)
}

// processSequential processes URLs one by one
func processSequential(urls []string) []Result {
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(10 * time.Second),
	)

	results := make([]Result, 0, len(urls))

	for _, url := range urls {
		start := time.Now()
		metadata, err := client.Extract(url)
		duration := time.Since(start)

		results = append(results, Result{
			URL:      url,
			Metadata: metadata,
			Error:    err,
			Duration: duration,
		})

		if err != nil {
			fmt.Printf("❌ %s - Error: %v\n", url, err)
		} else {
			fmt.Printf("✓ %s - %v\n", url, duration)
		}
	}

	return results
}

// processConcurrent processes URLs concurrently with worker pool
func processConcurrent(urls []string, workers int) []Result {
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(10 * time.Second),
	)

	// Channels for work distribution
	urlChan := make(chan string, len(urls))
	resultChan := make(chan Result, len(urls))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(i+1, client, urlChan, resultChan, &wg)
	}

	// Send URLs to workers
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := make([]Result, 0, len(urls))
	for result := range resultChan {
		results = append(results, result)
		if result.Error != nil {
			fmt.Printf("❌ Worker processed: %s - Error: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("✓ Worker processed: %s - %v\n", result.URL, result.Duration)
		}
	}

	return results
}

// worker processes URLs from the channel
func worker(id int, client *urlmeta.Client, urls <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range urls {
		start := time.Now()
		metadata, err := client.Extract(url)
		duration := time.Since(start)

		results <- Result{
			URL:      url,
			Metadata: metadata,
			Error:    err,
			Duration: duration,
		}
	}
}

// displayResults shows a summary of the results
func displayResults(results []Result) {
	successCount := 0
	errorCount := 0
	totalDuration := time.Duration(0)

	for _, result := range results {
		totalDuration += result.Duration
		if result.Error == nil {
			successCount++
		} else {
			errorCount++
		}
	}

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Total URLs: %d\n", len(results))
	fmt.Printf("Successful: %d\n", successCount)
	fmt.Printf("Failed: %d\n", errorCount)
	if len(results) > 0 {
		avgDuration := totalDuration / time.Duration(len(results))
		fmt.Printf("Average time per URL: %v\n", avgDuration)
	}
}

// Example of custom error handling and retry logic
func processWithRetry(url string, maxRetries int) (*urlmeta.Metadata, error) {
	client := urlmeta.NewClient(
		urlmeta.WithTimeout(10 * time.Second),
	)

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			fmt.Printf("Retry attempt %d for %s after %v\n", attempt, url, backoff)
			time.Sleep(backoff)
		}

		metadata, err := client.Extract(url)
		if err == nil {
			return metadata, nil
		}

		lastErr = err
		log.Printf("Attempt %d failed for %s: %v", attempt+1, url, err)
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}