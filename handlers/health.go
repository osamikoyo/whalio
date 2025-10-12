package handlers

import (
	"encoding/json"
	"net/http"
)

// Health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
	})
}
