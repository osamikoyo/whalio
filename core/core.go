package core

import (
	"context"
	"io"
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

func (c *Core) CreateArtist(name string, imagesource io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	artist := models.NewArtist(name)

	artist.ImagePath = artist.GetImageFilepath(c.cfg.ImageDir)

	if err := c.repository.CreateArtist(ctx, artist); err != nil {
		return err
	}

	return c.storage.SaveFile(imagesource, artist.ImagePath)
}

func (c *Core) PlaySong(id uint) (io.Reader, int64, error) {
	ctx, cancel := c.context()
	defer cancel()

	song, err := c.repository.GetSongByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	file, size, err := c.storage.OpenFile(song.Filepath(c.cfg.UploadDir))
	if err != nil {
		return nil, 0, err
	}

	return file, size, nil
}

func (c *Core) AddAlbum(name string, artistID uint, desc string, imagesource io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	album := models.NewAlbum(name, desc, artistID)

	album.ImagePath = album.ImageFilepath(c.cfg.ImageDir)

	if err := c.storage.SaveFile(imagesource, album.ImageFilepath(c.cfg.ImageDir)); err != nil {
		return err
	}

	if err := c.repository.CreateAlbum(ctx, album); err != nil {
		return err
	}

	return nil
}

func (c *Core) AddSong(name string, albumID uint, source io.Reader) error {
	ctx, cancel := c.context()
	defer cancel()

	song := models.NewSong(name, albumID)

	if err := c.storage.SaveFile(source, song.Filepath(c.cfg.UploadDir)); err != nil {
		return err
	}

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

