package handlers

import (
	"net/http"
	"whalio/templates"
)

// UploadSongsPage renders the songs upload page
func (h *Handlers) UploadSongsPage(w http.ResponseWriter, r *http.Request) {
	// Get all albums for the dropdown
	albums, err := h.core.GetSomeAlbums()
	if err != nil {
		h.SendError(w, r, "Failed to load albums", http.StatusInternalServerError)
		return
	}

	// Render the upload page with albums
	component := templates.UploadSongsPage(albums)
	if err := component.Render(r.Context(), w); err != nil {
		h.SendError(w, r, "Failed to render upload songs page", http.StatusInternalServerError)
		return
	}
}
