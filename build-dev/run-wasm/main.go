package main

import (
	"log"
	"net/http"
)

func main() {
	// Specify the directory containing your static files
	fs := http.FileServer(http.Dir("."))

	// Handle all requests with the file server
	http.Handle("/", fs)

	log.Print("Listening on port :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
