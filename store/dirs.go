package store

import (
	"os"
	"path"
)

import "crap/kvmap"

type Store struct {
	path string
	perm os.FileMode
}

func New(config kvmap.KVMap) *Store {
	path, err := config.GetString("store.path")
	if err != nil {
		panic(err)
	}
	perm, err := config.GetInt("store.permissions")
	if err != nil {
		panic(err)
	}

	s := Store{path, os.FileMode(perm)}
	s.createDirs()

	return &s
}

func (s Store) createDirs() {
	if err := os.MkdirAll(s.indexPath(), s.perm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(s.blobPath(), s.perm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(s.tempPath(), s.perm); err != nil {
		panic(err)
	}

	if err := cleanDir(s.tempPath()); err != nil {
		panic(err)
	}
	if err := syncDir(s.crapPath()); err != nil {
		panic(err)
	}
}

func (s Store) crapPath() string {
	return path.Join(s.path, "crap")
}

func (s Store) indexPath() string {
	return path.Join(s.crapPath(), "index")
}

func (s Store) blobPath() string {
	return path.Join(s.crapPath(), "blobs")
}

func (s Store) tempPath() string {
	return path.Join(s.crapPath(), "tmp")
}

func cleanDir(name string) error {
	dir, err := os.Open(name)
	if err != nil {
		return err
	}
	defer dir.Close()

	temps, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, t := range temps {
		if err := os.Remove(path.Join(name, t)); err != nil {
			return err
		}
	}

	return nil
}

func syncDir(name string) error {
	dir, err := os.Open(name)
	if err != nil {
		return err
	}
	defer dir.Close()

	if err := dir.Sync(); err != nil {
		return err
	}

	return nil
}
