package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"whalio/models"
	"whalio/templates"
)

// SearchContent handles search requests for albums, artists, and songs
func (h *Handlers) SearchContent(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		// If no query, return empty results
		if IsHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<div class="text-center text-base-content/50 py-4">Enter a search term to find music</div>`)
		} else {
			h.SendJSON(w, map[string]interface{}{"results": []interface{}{}}, http.StatusOK)
		}
		return
	}

	// Get all albums and artists
	albums, err := h.core.GetSomeAlbums()
	if err != nil {
		h.SendError(w, r, "Failed to search albums", http.StatusInternalServerError)
		return
	}

	artists, err := h.core.GetSomeArtist()
	if err != nil {
		h.SendError(w, r, "Failed to search artists", http.StatusInternalServerError)
		return
	}

	// Filter based on query
	var matchedAlbums []models.Album
	var matchedArtists []models.Artist

	queryLower := strings.ToLower(query)

	// Search in albums
	for _, album := range albums {
		if strings.Contains(strings.ToLower(album.Name), queryLower) ||
			strings.Contains(strings.ToLower(album.Description), queryLower) ||
			strings.Contains(strings.ToLower(album.Artist.Name), queryLower) {
			matchedAlbums = append(matchedAlbums, album)
		}
	}

	// Search in artists
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), queryLower) ||
			strings.Contains(strings.ToLower(artist.Desc), queryLower) {
			matchedArtists = append(matchedArtists, artist)
		}
	}

	if IsHTMXRequest(r) {
		// Return HTML for HTMX requests
		component := templates.SearchResults(matchedAlbums, matchedArtists, query)
		if err := component.Render(r.Context(), w); err != nil {
			h.SendError(w, r, "Failed to render search results", http.StatusInternalServerError)
			return
		}
	} else {
		// Return JSON for API requests
		result := map[string]interface{}{
			"query":   query,
			"albums":  matchedAlbums,
			"artists": matchedArtists,
		}
		h.SendJSON(w, result, http.StatusOK)
	}
}
