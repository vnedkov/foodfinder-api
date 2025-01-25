package handlers

import "net/http"

// HealthHandler is an HTTP handler that returns the health status of the API
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the response body
	w.Write([]byte(`{"status": "OK"}`))
}
