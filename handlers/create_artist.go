package handlers

import (
	"net/http"
)

// CreateArtist handles creating a new artist
func (h *Handlers) CreateArtist(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	desc := r.FormValue("desc")

	file, _, err := r.FormFile("file")
	if err != nil {
		h.SendError(w, r, "failed get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.core.CreateArtist(name, desc, file); err != nil {
		h.SendError(w, r, "failed create artist", http.StatusInternalServerError)
		return
	}

	if IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<div class="alert alert-success"><span>Artist created successfully</span></div>`))
		return
	}

	h.SendJSON(w, map[string]any{
		"success": true,
		"message": "Artist created successfully",
	}, http.StatusOK)
}
