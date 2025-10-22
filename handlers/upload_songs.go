package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// UploadSongs handles song file uploads
func (h *Handlers) UploadSongs(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.SendError(w, r, "Failed to parse form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	albumIDStr := r.FormValue("album_id")
	songTitle := strings.TrimSpace(r.FormValue("song_title"))

	// Validate album ID
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		h.SendError(w, r, "Invalid album ID", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, fileHeader, err := r.FormFile("audio_file")
	if err != nil {
		h.SendError(w, r, "No audio file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file
	if err := h.validateAudioFile(fileHeader); err != nil {
		h.SendError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Use filename as title if no title provided
	if songTitle == "" {
		songTitle = strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	}

	// Detect MIME type
	mimeType := h.detectMimeType(fileHeader.Filename)

	// Add song to database and save file
	if err := h.core.AddSong(songTitle, fileHeader.Filename, mimeType, fileHeader.Size, uint(albumID), file); err != nil {
		h.SendError(w, r, "Failed to upload song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if IsHTMXRequest(r) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<div class="alert alert-success"><span>✓ Song uploaded successfully</span></div>`)
	} else {
		h.SendJSON(w, map[string]interface{}{
			"success": true,
			"message": "Song uploaded successfully",
		}, http.StatusOK)
	}
}

// validateAudioFile checks if the uploaded file is a valid audio file
func (h *Handlers) validateAudioFile(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 100 << 20 // 100MB

	// Check file size
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file too large (max 100MB)")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	validExts := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".flac": true,
		".ogg":  true,
		".m4a":  true,
		".aac":  true,
	}

	if !validExts[ext] {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	return nil
}

// detectMimeType returns the MIME type based on file extension
func (h *Handlers) detectMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := map[string]string{
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".ogg":  "audio/ogg",
		".m4a":  "audio/m4a",
		".aac":  "audio/aac",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}

	return "audio/mpeg" // default
}
