package main

import (
	"log"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content- Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", readinessHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileserver))

	log.Println("Starting server on on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
