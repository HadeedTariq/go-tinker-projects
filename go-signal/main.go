package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Simulated database connection
type Database struct {
	mu        sync.Mutex
	connected bool
}

// Connect to the database
func (db *Database) Connect() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.connected = true
	fmt.Println("✅ Database connected")
}

// Disconnect the database
func (db *Database) Disconnect() {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.connected {
		fmt.Println("⏳ Closing database connection...")
		time.Sleep(2 * time.Second) // Simulating cleanup delay
		db.connected = false
		fmt.Println("✅ Database disconnected")
	}
}

func main() {
	// Create a simulated database
	db := &Database{}
	db.Connect() // Connect on startup

	// Create an HTTP server
	server := &http.Server{Addr: ":3000", Handler: http.DefaultServeMux}

	// Define a simple endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		fmt.Println("🚀 Server is running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for an OS signal (Ctrl+C or SIGTERM)
	sig := <-sigChan
	fmt.Println("\n🛑 Received signal:", sig)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("⏳ Shutting down server gracefully...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Close database connection
	db.Disconnect()

	fmt.Println("✅ Server shut down cleanly. Goodbye!")
}
