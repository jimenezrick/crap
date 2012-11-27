package store

// XXX: create(O_EXCLUSIVE) de file that is going to be renamed

// XXX: Encrypt AES
// XXX: Conn SSL
// XXX: Retrieve(key string) *Blob

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

import "crap/config"

type blob struct {
	file   *os.File
	writer *bufio.Writer
	hash   hash.Hash
}

func NewBlob() (*blob, error) {
	file, err := ioutil.TempFile(tempPath(), "blob")
	if err != nil {
		return nil, err
	}

	b := blob{file, bufio.NewWriter(file), sha1.New()}
	runtime.SetFinalizer(&b, func(b *blob) {
		os.Remove(b.file.Name())
	})

	return &b, nil
}

func (b *blob) Write(buf []byte) (int, error) {
	b.hash.Write(buf)
	return b.writer.Write(buf)
}

func (b *blob) Store() (string, error) {
	defer b.file.Close()
	b.writer.Flush()

	if err := b.file.Sync(); err != nil {
		return "", err
	}

	src := b.file.Name()
	dest := b.Path()
	perm := os.FileMode(config.GetInt("store.permissions"))

	if err := os.MkdirAll(path.Dir(dest), perm); err != nil {
		return "", err
	}

	if err := os.Rename(src, dest); err != nil {
		return "", err
	}

	if err := syncDir(path.Dir(dest)); err != nil {
		return "", err
	}

	runtime.SetFinalizer(b, nil)
	return b.Key(), nil
}

func (b *blob) Abort() error {
	if err := b.file.Close(); err != nil {
		return err
	}

	if err := os.Remove(b.file.Name()); err != nil {
		return err
	}

	runtime.SetFinalizer(b, nil)
	return nil
}

func (b *blob) Size() (int64, error) {
	info, err := b.file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (b *blob) Key() string {
	return fmt.Sprintf("%x", b.hash.Sum(nil))
}

func (b *blob) Path() string {
	hash := b.Key()
	return path.Join(blobPath(), hash[:2], hash[2:])
}
