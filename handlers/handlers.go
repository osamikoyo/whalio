package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"whalio/core"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	core *core.Core
}

func New(core *core.Core) *Handlers {
	return &Handlers{
		core: core,
	}
}

// RegisterRoutes registers all application routes
func (h *Handlers) RegisterRoutes(r *chi.Mux) {
	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))

	// Page routes
	r.Get("/", h.Index)
	r.Get("/about", h.About)
	r.Get("/album/{id}", h.Album)
	r.Get("/artist/{id}", h.Artist)
	r.Get("/library", h.Library)
	r.Get("/create/album", h.CreateAlbumPage)
	r.Get("/create/artist", h.CreateArtistPage)
	r.Get("/upload", h.UploadSongsPage)

	// Streaming routes
	r.Get("/stream/{id}", h.StreamAudio)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/create/album", h.CreateAlbum)
		r.Post("/create/artist", h.CreateArtist)
		r.Post("/songs/upload", h.UploadSongs)
		r.Get("/stats", h.GetStats)
		r.Get("/search", h.SearchContent)
		r.Get("/delete/album/{id}", h.DeleteAlbum)
		r.Get("/delete/artist/{id}", h.DeleteArtist)
		// Player endpoints
		r.Get("/song/{id}", h.GetSongInfo)
		// Album endpoints
		r.Get("/album/{id}/songs", h.GetAlbumSongs)
	})

	// Health check
	r.Get("/health", h.Health)
}

// Utility function to check if request is from HTMX
func IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// Utility function to get client IP
func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// Utility function to send JSON response
func (h *Handlers) SendJSON(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Utility function to send error response
func (h *Handlers) SendError(w http.ResponseWriter, r *http.Request, message string, status int) {
	if IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(status)
		alertClass := "alert-error"
		if status >= 400 && status < 500 {
			alertClass = "alert-warning"
		}

		fmt.Fprintf(w, `
			<div class="alert %s">
				<svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>%s</span>
			</div>
		`, alertClass, message)
		return
	}

	h.SendJSON(w, map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  status,
	}, status)
}

// ParseIntParam safely parses an integer parameter from URL
func ParseIntParam(r *http.Request, param string, defaultValue int) int {
	str := chi.URLParam(r, param)
	if str == "" {
		return defaultValue
	}

	if val, err := strconv.Atoi(str); err == nil {
		return val
	}

	return defaultValue
}
