package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.SendError(w, r, "failed convert id to int", http.StatusBadRequest)
		return
	}

	if err := h.core.DeleteAlbum(uint(id)); err != nil {
		h.SendError(w, r, "failed delete album", http.StatusInternalServerError)
		return
	}
}
