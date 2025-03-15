package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Log structure
var logMessages = []string{
	"INFO: Service started",
	"ERROR: Database connection failed",
	"INFO: Request handled",
	"ERROR: Disk space low",
	"INFO: User logged in",
	"ERROR: API request timeout",
}

// Worker function that processes logs and filters out error messages
func worker(id int, logs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for log := range logs {
		fmt.Printf("Worker %d processing log: %s\n", id, log)
		if strings.Contains(log, "ERROR") {
			results <- log // Send only error logs to results channel
		}
		time.Sleep(time.Millisecond * 500) // Simulate processing time
	}
}

func main() {
	logChannel := make(chan string, len(logMessages)) // Buffered channel for logs
	results := make(chan string, len(logMessages))    // Channel for filtered error logs
	var wg sync.WaitGroup

	// Start worker pool (3 workers)
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, logChannel, results, &wg)
	}

	// Send logs to logChannel
	go func() {
		for _, log := range logMessages {
			logChannel <- log
		}
		close(logChannel) // Close channel after sending all logs
	}()

	// Wait for workers to finish and then close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Read results (error logs)
	for res := range results {
		fmt.Println("Filtered Error Log:", res)
	}
}
