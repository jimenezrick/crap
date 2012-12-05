package store

import (
	"os"
	"path"
)

import "crap/util"

func (s Store) initStore() error {
	if err := os.MkdirAll(s.indexPath(), s.dirPerm); err != nil {
		return err
	}
	if err := os.MkdirAll(s.blobPath(), s.dirPerm); err != nil {
		return err
	}

	if exist, err := util.FileExist(s.tempPath()); err != nil {
		return err
	} else if exist {
		if err := os.RemoveAll(s.tempPath()); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(s.tempPath(), s.dirPerm); err != nil {
		return err
	}

	if err := util.Datasync(s.crapPath()); err != nil {
		return err
	}

	return nil
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

func (s Store) lockPath() string {
	return path.Join(s.crapPath(), "lock")
}
