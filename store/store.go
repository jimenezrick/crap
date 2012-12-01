package store

import "os"

import (
	"crap/config"
	"crap/util"
)

type Store struct {
	path     string
	dirPerm  os.FileMode
	filePerm os.FileMode
}

func Open(config config.Config) (*Store, error) {
	path := config.GetString("store.path")
	dirPerm := config.GetIntString("store.dir_permissions")
	filePerm := config.GetIntString("store.file_permissions")
	s := Store{path, os.FileMode(dirPerm), os.FileMode(filePerm)}

	if err := os.MkdirAll(s.crapPath(), s.dirPerm); err != nil {
		return nil, err
	}
	if err := util.CreateLockFile(s.lockPath(), s.filePerm); err != nil {
		return nil, err
	}

	if err := s.initStore(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s Store) Close() error {
	return os.Remove(s.lockPath())
}
