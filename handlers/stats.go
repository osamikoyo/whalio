package handlers

import (
	"net/http"
)

// StatsResponse represents the statistics data
type StatsResponse struct {
	Albums  int `json:"albums"`
	Artists int `json:"artists"`
	Songs   int `json:"songs"`
}

// GetStats returns statistics about the music library
func (h *Handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	// Get albums count
	albums, err := h.core.GetSomeAlbums()
	if err != nil {
		h.SendError(w, r, "Failed to get albums count", http.StatusInternalServerError)
		return
	}

	// Get artists count
	artists, err := h.core.GetSomeArtist()
	if err != nil {
		h.SendError(w, r, "Failed to get artists count", http.StatusInternalServerError)
		return
	}

	// Count total songs across all albums
	totalSongs := 0
	for _, album := range albums {
		totalSongs += len(album.Songs)
	}

	stats := StatsResponse{
		Albums:  len(albums),
		Artists: len(artists),
		Songs:   totalSongs,
	}

	h.SendJSON(w, stats, http.StatusOK)
}
