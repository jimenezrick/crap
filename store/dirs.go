package store

import (
	"os"
	"path"
)

import "crap/config"

func Init() error {
	perm := os.FileMode(config.GetInt("store.permissions"))

	if err := os.MkdirAll(indexPath(), perm); err != nil {
		return err
	}
	if err := os.MkdirAll(blobPath(), perm); err != nil {
		return err
	}
	if err := os.MkdirAll(tempPath(), perm); err != nil {
		return err
	}

	if err := cleanTempDir(); err != nil {
		return err
	}

	return nil
}

func crapPath() string {
	return path.Join(config.GetString("store.path"), "crap")
}

func indexPath() string {
	return path.Join(crapPath(), "index")
}

func blobPath() string {
	return path.Join(crapPath(), "blobs")
}

func tempPath() string {
	return path.Join(crapPath(), "tmp")
}

func cleanTempDir() error {
	dir, err := os.Open(tempPath())
	if err != nil {
		return err
	}
	defer dir.Close()

	temps, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, t := range temps {
		if err := os.Remove(path.Join(tempPath(), t)); err != nil {
			return err
		}
	}

	return nil
}
