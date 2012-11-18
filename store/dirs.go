package store

import (
	"os"
	"path"
)

const (
	pathStore = "crap"
	pathIndex = pathStore + "/" + "index"
	pathBlobs = pathStore + "/" + "blobs"
	pathTemp  = pathStore + "/" + "tmp"

	defaultPerm = 0700
)

func Init() error {
	if err := os.MkdirAll(pathStore, defaultPerm); err != nil {
		return err
	}

	if err := os.MkdirAll(pathIndex, defaultPerm); err != nil {
		return err
	}

	if err := os.MkdirAll(pathBlobs, defaultPerm); err != nil {
		return err
	}

	if err := os.MkdirAll(pathTemp, defaultPerm); err != nil {
		return err
	}

	dir, err := os.Open(pathTemp)
	if err != nil {
		return err
	}
	defer dir.Close()

	temps, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, t := range temps {
		if err := os.Remove(path.Join(pathTemp, t)); err != nil {
			return err
		}
	}

	return nil
}
