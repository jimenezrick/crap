package store

import "os"

import "crap/kvmap"

type Store struct {
	path string
	perm os.FileMode
}

func New(config *kvmap.KVMap) (*Store, error) {
	path, err := config.GetString("store.path")
	if err != nil {
		panic(err)
	}
	perm, err := config.GetInt("store.permissions")
	if err != nil {
		panic(err)
	}
	s := Store{path, os.FileMode(perm)}

	if err := s.createDirs(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s Store) Lock() error {
	file, err := os.OpenFile(
		s.lockPath(),
		os.O_WRONLY | os.O_CREATE | os.O_EXCL,
		s.perm)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func (s Store) Unlock() error {
	return os.Remove(s.lockPath())
}
