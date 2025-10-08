package repository

import (
	"context"
	"whalio/models"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var (
	ErrArtistNotFound = errors.New("artist not found")
	ErrAlbumNotFound  = errors.New("album not found")
	ErrSongNotFound   = errors.New("song not found")
)

type Repository struct {
	logger *zerolog.Logger
	db     *gorm.DB
}

func NewRepository(logger *zerolog.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}

func (r *Repository) CreateArtist(ctx context.Context, artist *models.Artist) error {
	log := r.logger.With().Str("method", "CreateArtist").Logger()
	log.Info().Str("name", artist.Name).Msg("Creating new artist")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(artist).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create artist")
		return errors.Wrap(err, "failed to create artist")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit artist creation")
	}
	log.Debug().Uint("id", artist.ID).Msg("Artist created successfully")
	return nil
}

func (r *Repository) GetArtistByID(ctx context.Context, id uint) (*models.Artist, error) {
	log := r.logger.With().Str("method", "GetArtistByID").Uint("id", id).Logger()
	log.Info().Msg("Fetching artist")

	var artist models.Artist
	err := r.db.WithContext(ctx).
		Preload("Albums.Songs").
		First(&artist, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Artist not found")
			return nil, ErrArtistNotFound
		}
		log.Error().Stack().Err(err).Msg("Failed to get artist")
		return nil, errors.Wrap(err, "failed to get artist")
	}
	log.Debug().Msg("Artist fetched successfully")
	return &artist, nil
}

func (r *Repository) GetArtistByName(ctx context.Context, name string) (*models.Artist, error) {
	log := r.logger.With().Str("method", "GetArtistByName").Str("name", name).Logger()
	log.Info().Msg("Fetching artist")

	var artist models.Artist
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&artist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Artist not found")
			return nil, ErrArtistNotFound
		}
		log.Error().Stack().Err(err).Msg("Failed to get artist")
		return nil, errors.Wrap(err, "failed to get artist")
	}
	log.Debug().Msg("Artist fetched successfully")
	return &artist, nil
}

func (r *Repository) ListArtists(ctx context.Context) ([]models.Artist, error) {
	log := r.logger.With().Str("method", "ListArtists").Logger()
	log.Info().Msg("Fetching artists")

	var albums []models.Artist
	err := r.db.WithContext(ctx).
		Find(&albums).Error
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to fetch artists")
		return nil, errors.Wrap(err, "failed to fetch artists")
	}

	if len(albums) == 0 {
		log.Warn().Msg("No artists found")
		return []models.Artist{}, nil
	}

	log.Debug().Int("count", len(albums)).Msg("Albums fetched successfully")
	return albums, nil
}

func (r *Repository) UpdateArtist(ctx context.Context, artist *models.Artist) error {
	log := r.logger.With().Str("method", "UpdateArtist").Uint("id", artist.ID).Logger()
	log.Info().Msg("Updating artist")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Save(artist).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to update artist")
		return errors.Wrap(err, "failed to update artist")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit artist update")
	}
	log.Debug().Msg("Artist updated successfully")
	return nil
}

func (r *Repository) DeleteArtist(ctx context.Context, id uint) error {
	log := r.logger.With().Str("method", "DeleteArtist").Uint("id", id).Logger()
	log.Info().Msg("Deleting artist")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Delete(&models.Artist{}, id).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to delete artist")
		return errors.Wrap(err, "failed to delete artist")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit artist deletion")
	}
	log.Debug().Msg("Artist deleted successfully")
	return nil
}

func (r *Repository) CreateAlbum(ctx context.Context, album *models.Album) error {
	log := r.logger.With().Str("method", "CreateAlbum").Str("name", album.Name).
		Uint("artist_id", album.ArtistID).
		Str("image_path", album.ImagePath).
		Logger()
	log.Info().Msg("Creating new album")

	var artist models.Artist
	if err := r.db.WithContext(ctx).First(&artist, album.ArtistID).Error; err != nil {
		log.Warn().Err(err).Msg("Artist not found for album")
		return errors.Wrap(ErrArtistNotFound, "invalid artist ID for album")
	}

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(album).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to create album")
		return errors.Wrap(err, "failed to create album")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit album creation")
	}
	log.Debug().Uint("id", album.ID).Msg("Album created successfully")
	return nil
}

func (r *Repository) GetAlbumByID(ctx context.Context, id uint) (*models.Album, error) {
	log := r.logger.With().Str("method", "GetAlbumByID").Uint("id", id).Logger()
	log.Info().Msg("Fetching album")

	var album models.Album
	err := r.db.WithContext(ctx).
		Preload("Songs").
		Preload("Artist").
		First(&album, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Album not found")
			return nil, ErrAlbumNotFound
		}
		log.Error().Stack().Err(err).Msg("Failed to get album")
		return nil, errors.Wrap(err, "failed to get album")
	}
	log.Debug().Msg("Album fetched successfully")
	return &album, nil
}

func (r *Repository) ListAlbums(ctx context.Context) ([]models.Album, error) {
	log := r.logger.With().Str("method", "ListAlbums").Logger()
	log.Info().Msg("Fetching albums")

	var albums []models.Album
	err := r.db.WithContext(ctx).
		Preload("Artist").
		Preload("Songs").
		Find(&albums).Error
	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to fetch albums")
		return nil, errors.Wrap(err, "failed to fetch albums")
	}

	if len(albums) == 0 {
		log.Warn().Msg("No albums found")
		return []models.Album{}, nil
	}

	log.Debug().Int("count", len(albums)).Msg("Albums fetched successfully")
	return albums, nil
}
func (r *Repository) UpdateAlbum(ctx context.Context, album *models.Album) error {
	log := r.logger.With().Str("method", "UpdateAlbum").Uint("id", album.ID).Logger()
	log.Info().Msg("Updating album")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Save(album).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to update album")
		return errors.Wrap(err, "failed to update album")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit album update")
	}
	log.Debug().Msg("Album updated successfully")
	return nil
}

func (r *Repository) DeleteAlbum(ctx context.Context, id uint) error {
	log := r.logger.With().Str("method", "DeleteAlbum").Uint("id", id).Logger()
	log.Info().Msg("Deleting album")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Delete(&models.Album{}, id).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to delete album")
		return errors.Wrap(err, "failed to delete album")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit album deletion")
	}
	log.Debug().Msg("Album deleted successfully")
	return nil
}

func (r *Repository) CreateSong(ctx context.Context, song *models.Song) error {
	log := r.logger.With().Str("method", "CreateSong").Str("name", song.Name).Uint("album_id", song.AlbumID).Logger()
	log.Info().Msg("Creating new song")

	var album models.Album
	if err := r.db.WithContext(ctx).First(&album, song.AlbumID).Error; err != nil {
		log.Warn().Err(err).Msg("Album not found for song")
		return errors.Wrap(ErrAlbumNotFound, "invalid album ID for song")
	}

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(song).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to create song")
		return errors.Wrap(err, "failed to create song")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit song creation")
	}
	log.Debug().Uint("id", song.ID).Msg("Song created successfully")
	return nil
}

func (r *Repository) GetSongByID(ctx context.Context, id uint) (*models.Song, error) {
	log := r.logger.With().Str("method", "GetSongByID").Uint("id", id).Logger()
	log.Info().Msg("Fetching song")

	var song models.Song
	err := r.db.WithContext(ctx).
		Preload("Album.Artist").
		First(&song, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Song not found")
			return nil, ErrSongNotFound
		}
		log.Error().Stack().Err(err).Msg("Failed to get song")
		return nil, errors.Wrap(err, "failed to get song")
	}
	log.Debug().Msg("Song fetched successfully")
	return &song, nil
}

func (r *Repository) UpdateSong(ctx context.Context, song *models.Song) error {
	log := r.logger.With().Str("method", "UpdateSong").Uint("id", song.ID).Logger()
	log.Info().Msg("Updating song")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Save(song).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to update song")
		return errors.Wrap(err, "failed to update song")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit song update")
	}
	log.Debug().Msg("Song updated successfully")
	return nil
}

func (r *Repository) DeleteSong(ctx context.Context, id uint) error {
	log := r.logger.With().Str("method", "DeleteSong").Uint("id", id).Logger()
	log.Info().Msg("Deleting song")

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Delete(&models.Song{}, id).Error; err != nil {
		tx.Rollback()
		log.Error().Stack().Err(err).Msg("Failed to delete song")
		return errors.Wrap(err, "failed to delete song")
	}
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return errors.Wrap(err, "failed to commit song deletion")
	}
	log.Debug().Msg("Song deleted successfully")
	return nil
}

func (r *Repository) AddSongToAlbum(ctx context.Context, albumID uint, song *models.Song) error {
	log := r.logger.With().Str("method", "AddSongToAlbum").Uint("album_id", albumID).Str("song_name", song.Name).Logger()
	log.Info().Msg("Adding song to album")

	var album models.Album
	if err := r.db.WithContext(ctx).First(&album, albumID).Error; err != nil {
		log.Warn().Err(err).Msg("Album not found")
		return errors.Wrap(ErrAlbumNotFound, "invalid album ID for song")
	}

	song.AlbumID = albumID
	return r.CreateSong(ctx, song)
}
