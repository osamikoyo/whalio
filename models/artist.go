package models

import (
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

const ArtistImageTemplate = "%s.png"

type Artist struct {
	gorm.Model
	Name      string
	ImagePath string
	Desc      string
	Albums    []Album `gorm:"foreignKey:ArtistID"`
}

func NewArtist(name string, desc string) *Artist {
	return &Artist{
		Name: name,
		Desc: desc,
	}
}

func (a *Artist) GetImageFilepath(imageDir string) string {
	return filepath.Join(imageDir, fmt.Sprintf(ArtistImageTemplate, a.Name))
}
