package handlers

import (
	"net/http"
	"whalio/templates"
)

// Index renders the main page
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	component := templates.Index()
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
