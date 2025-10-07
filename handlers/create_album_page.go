package handlers

import (
	"net/http"
	"whalio/templates"
)

func (h *Handlers) CreateAlbumPage(w http.ResponseWriter, r *http.Request) {
	component := templates.CreateAlbumPage()
	if err := component.Render(r.Context(), w); err != nil {
		h.SendError(w, r, "Failed to render create album page", http.StatusInternalServerError)
		return
	}
}
