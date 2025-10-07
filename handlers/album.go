package handlers

import (
	"net/http"
	"strconv"
	"whalio/templates"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) Album(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendError(w, r, "failed convert id param to int", http.StatusBadRequest)
		return
	}

	album, err := h.core.GetAlbum(uint(id))
	if err != nil {
		h.SendError(w, r, "failed get album", http.StatusInternalServerError)
		return
	}

	component := templates.Album(album)
	if err = component.Render(r.Context(), w); err != nil {
		h.SendError(w, r, "failed render component", http.StatusInternalServerError)
		return
	}
}

// GetAlbumSongs returns list of songs for album playback
func (h *Handlers) GetAlbumSongs(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendError(w, r, "invalid album id", http.StatusBadRequest)
		return
	}
	album, err := h.core.GetAlbum(uint(id))
	if err != nil {
		h.SendError(w, r, "failed get album", http.StatusInternalServerError)
		return
	}
	type songDTO struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		AlbumID  uint   `json:"albumId"`
		Artist   string `json:"artist"`
		MimeType string `json:"mimeType"`
	}
	songs := make([]songDTO, 0, len(album.Songs))
	for _, s := range album.Songs {
		songs = append(songs, songDTO{ID: s.ID, Name: s.Name, AlbumID: s.AlbumID, Artist: album.Artist.Name, MimeType: s.MimeType})
	}
	_ = h.SendJSON(w, map[string]any{
		"albumId": album.ID,
		"album":   album.Name,
		"artist":  album.Artist.Name,
		"songs":   songs,
	}, http.StatusOK)
}
