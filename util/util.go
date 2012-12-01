package util

import (
	"fmt"
	"hash"
	"os"
)

func SyncFile(name string) error {
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

func FileExist(name string) (bool, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func CreateLockFile(name string, perm os.FileMode) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func HexHash(hash hash.Hash) string {
	return fmt.Sprintf("%x", hash.Sum(nil))
}
