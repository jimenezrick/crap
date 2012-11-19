package store

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

type Blob struct {
	file   *os.File
	writer *bufio.Writer
	hash   hash.Hash
}

func NewBlob() (*Blob, error) {
	file, err := ioutil.TempFile(tempPath(), "blob")
	if err != nil {
		return nil, err
	}

	blob := Blob{file, bufio.NewWriter(file), sha1.New()}
	runtime.SetFinalizer(&blob, func(b *Blob) {
		os.Remove(b.file.Name())
	})

	return &blob, nil
}

func (b *Blob) Write(buf []byte) (int, error) {
	b.hash.Write(buf)
	return b.writer.Write(buf)
}

func (b *Blob) Store() (string, error) {
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

	SyncFile(path.Dir(dest))
	return b.Key(), nil
}

func (b *Blob) Abort() error {
	if err := b.file.Close(); err != nil {
		return err
	}

	if err := os.Remove(b.file.Name()); err != nil {
		return err
	}

	return nil
}

func (b *Blob) Size() (int64, error) {
	info, err := b.file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (b *Blob) Key() string {
	return fmt.Sprintf("%x", b.hash.Sum(nil))
}

func (b *Blob) Path() string {
	hash := b.Key()
	return path.Join(blobPath(), hash[:2], hash[2:])
}
