package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON sends any Go value back to the client as a JSON response
// with the given HTTP status code. Every handler uses this so that
// responses are always formatted the same way.
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeError sends a JSON error response, e.g. {"error": "something went wrong"}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
