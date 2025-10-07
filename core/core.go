package core

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"
	"whalio/config"
	"whalio/models"
	"whalio/repository"
	"whalio/storage"
)

type Core struct {
	repository *repository.Repository
	storage    *storage.Storage
	cfg        *config.Config
	timeout    time.Duration
}

func NewCore(repository *repository.Repository, storage *storage.Storage, cfg *config.Config, timeout time.Duration) *Core {
	return &Core{
		repository: repository,
		storage:    storage,
		cfg:        cfg,
		timeout:    timeout,
	}
}

func (c *Core) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *Core) CreateArtist(name, desc string, imagesource io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	artist := models.NewArtist(name, desc)

	artist.ImagePath = artist.GetImageFilepath(c.cfg.ImageDir)

	if err := c.repository.CreateArtist(ctx, artist); err != nil {
		return err
	}

	return c.storage.SaveFile(imagesource, artist.ImagePath)
}

func (c *Core) PlaySong(id uint) (io.ReadSeeker, os.FileInfo, error) {
	ctx, cancel := c.context()
	defer cancel()

	song, err := c.repository.GetSongByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	file, info, err := c.storage.OpenFile(song.Filepath(c.cfg.UploadDir))
	if err != nil {
		return nil, nil, err
	}

	return file, info, nil
}

func (c *Core) AddSong(name, filename, mimeType string, fileSize int64, albumID uint, source io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	song := models.NewSong(name, filename, mimeType, fileSize, albumID)

	// Get album info for filepath generation
	album, err := c.repository.GetAlbumByID(ctx, albumID)
	if err != nil {
		return err
	}
	song.Album = *album

	// Save file to storage
	filePath := song.Filepath(c.cfg.UploadDir)
	if err := c.storage.SaveFile(source, filePath); err != nil {
		return err
	}

	// Create song in database
	if err := c.repository.CreateSong(ctx, song); err != nil {
		return err
	}

	return nil
}

func (c *Core) ChangeAlbum(songID uint, albumID uint) error {
	ctx, cancel := c.context()
	defer cancel()

	song, err := c.repository.GetSongByID(ctx, songID)
	if err != nil {
		return err
	}

	path := song.Filepath(c.cfg.UploadDir)

	song.AlbumID = albumID

	if err = c.repository.UpdateSong(ctx, song); err != nil {
		return err
	}

	newpath := song.Filepath(c.cfg.UploadDir)

	return c.storage.RenameFile(newpath, path)
}

func (c *Core) GetSomeAlbums() ([]models.Album, error) {
	ctx, cancel := c.context()
	defer cancel()

	return c.repository.ListAlbums(ctx)
}

func (c *Core) GetSomeArtist() ([]models.Artist, error) {
	ctx, cancel := c.context()
	defer cancel()

	return c.repository.ListArtists(ctx)
}

// GetSongByID returns song by ID with album and artist preloaded
func (c *Core) GetSongByID(id uint) (*models.Song, error) {
	ctx, cancel := c.context()
	defer cancel()
	return c.repository.GetSongByID(ctx, id)
}

func (c *Core) GetAlbum(id uint) (*models.Album, error) {
	ctx, cancel := c.context()
	defer cancel()

	return c.repository.GetAlbumByID(ctx, id)
}

func (c *Core) GetArtist(id uint) (*models.Artist, error) {
	ctx, cancel := c.context()
	defer cancel()

	return c.repository.GetArtistByID(ctx, id)
}

func (c *Core) CreateAlbum(name, desc, artistName string, year int, imageSource io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	artist, err := c.repository.GetArtistByName(ctx, artistName)
	if err != nil {
		return err
	}

	album := models.NewAlbum(name, desc, year, artist.ID)

	album.ImagePath = album.ImageFilepath()

	if err = c.repository.CreateAlbum(ctx, album); err != nil {
		return err
	}

	if err = c.storage.SaveFile(imageSource, filepath.Join(c.cfg.ImageDir, album.ImagePath)); err != nil {
		return err
	}

	return nil
}

func (c *Core) DeleteAlbum(id uint) error {
	ctx, cancel := c.context()
	defer cancel()

	album, err := c.repository.GetAlbumByID(ctx, id)
	if err != nil {
		return err
	}

	if err = c.repository.DeleteAlbum(ctx, id); err != nil {
		return err
	}

	if err = c.storage.DeleteFile(filepath.Join(c.cfg.ImageDir, album.ImagePath)); err != nil {
		return err
	}

	return nil
}

func (c *Core) DeleteArtist(id uint) error {
	ctx, cancel := c.context()
	defer cancel()

	artist, err := c.repository.GetArtistByID(ctx, id)
	if err != nil {
		return err
	}

	if err = c.repository.DeleteArtist(ctx, id); err != nil {
		return err
	}

	if err = c.storage.DeleteFile(filepath.Join(c.cfg.ImageDir, artist.ImagePath));err != nil{
		return err
	}

	return nil
}
