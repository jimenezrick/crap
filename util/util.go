package util

import (
	"fmt"
	"hash"
	"os"
	"syscall"
)

const (
	FALLOC_FL_KEEP_SIZE = 1 + iota
	FALLOC_FL_PUNCH_HOLE
)

const (
	POSIX_FADV_NORMAL = iota
	POSIX_FADV_RANDOM
	POSIX_FADV_SEQUENTIAL
	POSIX_FADV_WILLNEED
	POSIX_FADV_DONTNEED
	POSIX_FADV_NOREUSE
)

func Fdatasync(file *os.File) error {
	return syscall.Fdatasync(int(file.Fd()))
}

func Datasync(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	return Fdatasync(file)
}

func Fallocate(file *os.File, mode uint32, off int64, len int64) error {
	return syscall.Fallocate(int(file.Fd()), mode, off, len)
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
