package storage

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Storage struct {
	logger *zerolog.Logger
}

func NewStorage(logger *zerolog.Logger) *Storage {
	return &Storage{
		logger: logger,
	}
}

func (s *Storage) SaveFile(source io.Reader, name string) error {
	dest, err := os.Create(name)
	if err != nil {
		s.logger.Error().Msgf("failed create dest for %s: %v", name, err)

		return err
	}
	defer dest.Close()

	if _, err = io.Copy(dest, source); err != nil {
		s.logger.Error().Msgf("failed copy %s: %v", name, err)

		return err
	}

	return nil
}

func (s *Storage) DeleteFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		s.logger.Error().Msgf("failed remove %s: %v", name, err)

		return err
	}

	return nil
}

func (s *Storage) ListFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		s.logger.Error().Msgf("failed read dir %s: %v", dir, err)

		return nil, err
	}

	names := make([]string, len(files))

	for i, file := range files {
		names[i] = file.Name()
	}

	return names, nil
}

func (s *Storage) RenameFile(destPath string, srcPath string) error {
	dest, err := os.Open(destPath)
	if err != nil{
		s.logger.Error().Msgf("failed open dest: %s: %v", destPath, err)
		return err
	}

	src, err := os.Open(srcPath)
	if err != nil{
		s.logger.Error().Msgf("failed open source: %s: %v", srcPath, err)
		return err
	}


	if _, err = io.Copy(dest, src);err != nil{
		s.logger.Error().Msgf("failed copy: %v", err)
		return err
	}

	return nil
}

func (s *Storage) GetFile(dest io.Writer, name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		s.logger.Error().Msgf("failed get stat of %s: %v", name, err)
		return err
	}

	source, err := os.Open(name)
	if err != nil {
		s.logger.Error().Msgf("failed open: %s: %v", name, err)

		return err
	}

	if _, err = io.Copy(dest, source); err != nil {
		s.logger.Error().Msgf("failed copy %s: %v", name, err)

		return err
	}

	return nil
}

func (s *Storage) OpenFile(name string) (io.Reader, int64, error) {
	stat, err := os.Stat(name)
	if os.IsNotExist(err) {
		s.logger.Error().Msgf("failed get stat of %s: %v", name, err)
		return nil, 0, err
	}

	source, err := os.Open(name)
	if err != nil {
		s.logger.Error().Msgf("failed open: %s: %v", name, err)

		return nil, 0, err
	}

	return source, stat.Size(), nil
}