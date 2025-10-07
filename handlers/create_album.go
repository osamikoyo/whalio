package handlers

import (
	"net/http"
	"strconv"
)

func (h *Handlers) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	yearStr := r.FormValue("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.SendError(w, r, "failed convert year to int", http.StatusBadRequest)
		return
	}

	artist := r.FormValue("artist")
	desc := r.FormValue("desc")

	file, _, err := r.FormFile("file")
	if err != nil {
		h.SendError(w, r, "failed get file from form", http.StatusBadRequest)
		return
	}

	if err = h.core.CreateAlbum(name, desc, artist, year, file); err != nil {
		h.SendError(w, r, "failed create album", http.StatusInternalServerError)
		return
	}
}
