package handlers

import (
	"net/http"
	"strconv"
	"whalio/templates"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) Artist(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		h.SendError(w, r, "failed convert id to int", http.StatusBadRequest)
		return
	}

	artist, err := h.core.GetArtist(uint(id))
	if err != nil {
		h.SendError(w, r, "failed get artist", http.StatusInternalServerError)
		return
	}

	component := templates.Artist(artist)
	if err = component.Render(r.Context(), w); err != nil {
		h.SendError(w, r, "failed render page", http.StatusInternalServerError)
		return
	}
}
