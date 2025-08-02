// ----------------------------------------------------
// FILE: backend/api.go
// The HTTP Server and request handlers.
// ----------------------------------------------------
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIServer struct remains the same...
type APIServer struct {
	listenAddr string
	store      Store
}

func NewAPIServer(listenAddr string, store Store) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Start() error {
	// Using http.ServeMux for better routing
	mux := http.NewServeMux()
	// NEW: Add a health check handler for the root path.
	mux.HandleFunc("/", s.handleHealthCheck) 
	mux.HandleFunc("/profiles/", s.handleGetProfile)
	
	log.Println("API Server starting on", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, mux)
}

// NEW: Health check handler function.
func (s *APIServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    // If the path is not exactly "/", treat it as a 404.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *APIServer) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	slug := r.URL.Path[len("/profiles/"):] // Simple routing
	if slug == "" {
		http.Error(w, "Profile slug is required", http.StatusBadRequest)
		return
	}
	
	profile, err := s.store.GetProfileBySlug(slug)
	if err != nil {
		// In a real app, differentiate between 404 Not Found and 500 Server Error
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	
	s.writeJSON(w, http.StatusOK, profile)
}

// writeJSON is a helper function for sending JSON responses.
func (s *APIServer) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}