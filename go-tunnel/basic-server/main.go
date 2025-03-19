package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler) // Handle all requests at "/"

	fmt.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil) // Start server on port 3000
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
