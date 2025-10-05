package models

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

const PathKeyString = "%s-%s-%s.mp3"

type Song struct {
    gorm.Model
    Name    string
    AlbumID uint 
    Album   Album  `gorm:"foreignKey:AlbumID"` 
}

func NewSong(name string, albumID uint) *Song {
	return &Song{
		Name: name,
		AlbumID: albumID,
	}
}

func (s *Song) Filepath(uploadDir string) string {
	return filepath.Join(uploadDir, fmt.Sprintf(PathKeyString, s.Name, s.Album.Name, s.Album.Artist.Name))
}

