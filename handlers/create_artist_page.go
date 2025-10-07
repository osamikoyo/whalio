package handlers

import (
	"net/http"
	"whalio/templates"
)

func (h *Handlers) CreateArtistPage(w http.ResponseWriter, r *http.Request) {
	component := templates.CreateArtistPage()
	if err := component.Render(r.Context(), w); err != nil {
		h.SendError(w, r, "failed render page", http.StatusBadRequest)
		return
	}
}
