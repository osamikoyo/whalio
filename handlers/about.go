package handlers

import (
	"net/http"
	"whalio/templates"
)

// About renders the about page
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	component := templates.About()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}