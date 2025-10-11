package handlers

import (
	"net/http"
	"whalio/templates"
)

func (h *Handlers) Library(w http.ResponseWriter, r *http.Request) {
	albums, err := h.core.GetSomeAlbums()
	if err != nil {
		h.SendError(w, r, "failed load albums", http.StatusInternalServerError)
		return
	}

	artists, err := h.core.GetSomeArtist()
	if err != nil {
		h.SendError(w, r, "failed load artists", http.StatusInternalServerError)
		return
	}

	if err := templates.Library(albums, artists).Render(r.Context(), w); err != nil {
		h.SendError(w, r, "failed render", http.StatusInternalServerError)
	}
}
