package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT") // Set different ports for each server
	if port == "" {
		port = "5001"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on port %s\n", port)
	})

	fmt.Printf("Server started at :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
