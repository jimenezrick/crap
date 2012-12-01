package store

import "os"

import "crap/config"

type Store struct {
	path string
	dirPerm os.FileMode
	filePerm os.FileMode
}

func New(config config.Config) (*Store, error) {
	path := config.GetString("store.path")
	dirPerm := config.GetIntString("store.dir_permissions")
	filePerm := config.GetIntString("store.file_permissions")

	s := Store{path, os.FileMode(dirPerm), os.FileMode(filePerm)}
	if err := s.initStore(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s Store) Lock() error {
	file, err := os.OpenFile(
		s.lockPath(),
		os.O_WRONLY | os.O_CREATE | os.O_EXCL,
		s.filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (s Store) Unlock() error {
	return os.Remove(s.lockPath())
}
