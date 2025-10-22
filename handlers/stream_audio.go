package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// StreamAudio handles streaming audio files with range request support
func (h *Handlers) StreamAudio(w http.ResponseWriter, r *http.Request) {
	// Get song ID from URL parameter
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		h.SendError(w, r, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Get song file from core
	file, fileInfo, err := h.core.PlaySong(uint(songID))
	if err != nil {
		h.SendError(w, r, "Song not found", http.StatusNotFound)
		return
	}
	// Close file if it supports Close()
	if closer, ok := file.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Set content type and length
	w.Header().Set("Content-Type", "audio/mpeg") // Default to MP3, could be improved to detect actual type
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileInfo.Name()))

	fileSize := fileInfo.Size()

	// Handle range requests for audio streaming
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		h.handleRangeRequest(w, r, file, fileSize, rangeHeader)
		return
	}

	// No range request, serve full file
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		// Don't send error response as we've already started writing
		return
	}
}

// handleRangeRequest handles HTTP range requests for partial content delivery
func (h *Handlers) handleRangeRequest(w http.ResponseWriter, r *http.Request, file io.ReadSeeker, fileSize int64, rangeHeader string) {
	// Parse range header (format: "bytes=start-end")
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	rangeValue := rangeHeader[6:] // Remove "bytes=" prefix
	ranges := strings.Split(rangeValue, ",")

	// Handle only the first range for simplicity
	if len(ranges) == 0 {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	rangeParts := strings.Split(ranges[0], "-")
	if len(rangeParts) != 2 {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	var start, end int64
	var err error

	// Parse start
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil || start < 0 {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	// Parse end
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil || end >= fileSize {
			end = fileSize - 1
		}
	} else {
		end = fileSize - 1
	}

	// Validate range
	if start > end || start >= fileSize {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Seek to start position
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		http.Error(w, "Failed to seek file", http.StatusInternalServerError)
		return
	}

	// Set headers for partial content
	contentLength := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	w.WriteHeader(http.StatusPartialContent)

	// Copy the requested range
	if _, err := io.CopyN(w, file, contentLength); err != nil {
		// Don't send error response as we've already started writing
		return
	}
}

// GetSongInfo returns information about a song for the player
func (h *Handlers) GetSongInfo(w http.ResponseWriter, r *http.Request) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		h.SendError(w, r, "Invalid song ID", http.StatusBadRequest)
		return
	}

	// Get song from repository
	song, err := h.core.GetSongByID(uint(songID))
	if err != nil {
		h.SendError(w, r, "Song not found", http.StatusNotFound)
		return
	}

	songInfo := map[string]interface{}{
		"id":       song.ID,
		"name":     song.Name,
		"filename": song.Filename,
		"mimeType": song.MimeType,
		"fileSize": song.FileSize,
		"duration": song.Duration,
		"album": map[string]interface{}{
			"id":   song.Album.ID,
			"name": song.Album.Name,
			"year": song.Album.Year,
			"artist": map[string]interface{}{
				"id":   song.Album.Artist.ID,
				"name": song.Album.Artist.Name,
			},
		},
	}

	h.SendJSON(w, songInfo, http.StatusOK)
}
