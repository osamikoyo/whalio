package models

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

const ImageFilepathKey = "%s:%d.png"

type Album struct {
	gorm.Model
	Name        string
	Description string
	ArtistID    uint
	Year int
	ImagePath string
	Artist      Artist `gorm:"foreignKey:ArtistID"`
	Songs       []Song `gorm:"foreignKey:AlbumID"`
}

func NewAlbum(name string, desc string, year int, artistID uint) *Album {
	return &Album{
		Name:        name,
		Description: desc,
		ArtistID:    artistID,
		Year: year,
	}
}

func (a *Album) ImageFilepath(imageDir string) string {
	return filepath.Join(imageDir, fmt.Sprintf(ImageFilepathKey, a.Name, a.ArtistID))
}

