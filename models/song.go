package models

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

const PathKeyString = "%s-%s-%s%s" // filename format: song-album-artist.ext

type Song struct {
	gorm.Model
	Name     string
	Filename string // Original filename with extension
	MimeType string // e.g., "audio/mpeg", "audio/wav"
	FileSize int64  // Size in bytes
	Duration int    // Duration in seconds
	AlbumID  uint
	Album    Album `gorm:"foreignKey:AlbumID"`
}

func NewSong(name, filename, mimeType string, fileSize int64, albumID uint) *Song {
	return &Song{
		Name:     name,
		Filename: filename,
		MimeType: mimeType,
		FileSize: fileSize,
		AlbumID:  albumID,
	}
}

func (s *Song) Filepath(uploadDir string) string {
	ext := filepath.Ext(s.Filename)
	return filepath.Join(uploadDir, fmt.Sprintf(PathKeyString, s.Name, s.Album.Name, s.Album.Artist.Name, ext))
}

// GetFileExtension returns the file extension from filename
func (s *Song) GetFileExtension() string {
	return filepath.Ext(s.Filename)
}

// IsAudioFile checks if the file is a valid audio format
func (s *Song) IsAudioFile() bool {
	validMimeTypes := map[string]bool{
		"audio/mpeg":   true, // MP3
		"audio/wav":    true, // WAV
		"audio/ogg":    true, // OGG
		"audio/m4a":    true, // M4A
		"audio/flac":   true, // FLAC
		"audio/x-flac": true, // FLAC alternative
	}
	return validMimeTypes[s.MimeType]
}
